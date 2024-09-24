<?php declare(strict_types = 1);

/*
 * This file is part of the Vairogs package.
 *
 * (c) Dāvis Zālītis (k0d3r1s) <davis@vairogs.com>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Vairogs\Component\CacheWarmer;

use RuntimeException;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;
use Symfony\Component\HttpClient\HttpClient;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Process\Process;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;
use Vairogs\Component\Functions\Local;
use Vairogs\Component\Functions\Php;
use Vairogs\Component\Functions\Web;

use function array_values;
use function chmod;
use function fclose;
use function file_put_contents;
use function fopen;
use function fwrite;
use function getcwd;
use function is_dir;
use function is_file;
use function is_resource;
use function is_writable;
use function sprintf;
use function stripos;

use const PHP_EOL;
use const PHP_OS_FAMILY;

final class CacheWarmerBinary
{
    public const string REPOSITORY = 'lettland/cache-warmer';
    private ?SymfonyStyle $output = null;

    public function __construct(
        private readonly ParameterBagInterface $bag,
        private ?HttpClientInterface $httpClient = null,
    ) {
        $this->httpClient ??= HttpClient::create();
    }

    /**
     * @throws TransportExceptionInterface
     */
    public function createProcess(
        array $arguments,
        bool $manual = false,
    ): Process {
        if (null === $this->output) {
            throw new RuntimeException(sprintf('%s can only be used with %s. Call ->setOutput($style) before %s function', __CLASS__, SymfonyStyle::class, __FUNCTION__));
        }

        static $_helper = null;

        if (null === $_helper) {
            $_helper = new class {
                use Local\_CurlUA;
                use Local\_MkDir;
                use Php\_SystemInfo;
                use Web\_LatestReleaseTag;
            };
        }

        $info = $_helper->systemInfo();

        if (!$manual) {
            $process = new Process(0 === stripos(PHP_OS_FAMILY, 'WIN') ? ['ping', '-n', '4', 'github.com'] : ['ping', '-c', '4', 'github.com']);
            $process->setTimeout(5)->run();

            if (!$process->isSuccessful()) {
                $version = 'nightly' === $this->bag->get('vairogs.cache_warmer.version') ? 'nightly' : 'latest';
                $url = sprintf('https://github.com/%s/releases/download/%s/cache-warmer-%s-%s%s', self::REPOSITORY, $version, ...array_values($info));

                throw new RuntimeException(sprintf('Unable to connect to github!%sPlease download desired binary from %s, put it in var/vairogs/bin/manual (name MUST match as is in github) and run this command again with --manual option', PHP_EOL, $url));
            }

            $configVersion = $this->bag->get('vairogs.cache_warmer.version');
            $version = match ($configVersion) {
                'latest' => $_helper->latestReleaseTag(self::REPOSITORY),
                default => $configVersion,
            };
        } else {
            $version = 'manual';
        }

        $directory = sprintf('%s/var/%s/%s', $this->bag->get('kernel.project_dir'), self::REPOSITORY, $version);

        $isDir = true;

        if (!is_dir($directory)) {
            $isDir = $_helper->mkdir($directory);
        }

        if (!$isDir || !is_writable($directory)) {
            throw new RuntimeException(sprintf('Unable to create or write directory %s', $directory));
        }

        if (!is_file($directory . '/.gitignore')) {
            file_put_contents($directory . '/.gitignore', '*' . PHP_EOL);
        }

        $path = sprintf('cache-warmer-%s-%s%s', ...array_values($info));
        $binary = $directory . '/' . $path;

        if (!$manual && !is_file($binary)) {
            $url = sprintf('https://github.com/%s/releases/download/%s/%s', self::REPOSITORY, $version, $path);

            $response = $this->httpClient->request(Request::METHOD_HEAD, $url);
            $statusCode = $response->getStatusCode();

            if (Response::HTTP_OK !== $statusCode) {
                throw new RuntimeException(sprintf('Unable to download binary from %s. Please check that version tag is set correctly in configuration', $url));
            }

            $this->output->note(sprintf('Downloading %s binary from %s', self::REPOSITORY, $url));
            $progressBar = null;
            $response = $this->httpClient->request(Request::METHOD_GET, $url, [
                'on_progress' => function (int $dlNow, int $dlSize) use (&$progressBar): void {
                    if (0 === $dlSize) {
                        return;
                    }

                    if (!$progressBar) {
                        $progressBar = $this->output->createProgressBar($dlSize);
                    }

                    $progressBar?->setProgress($dlNow);
                },
            ]);

            $fileHandler = fopen($binary, 'wb');

            if (!is_resource($fileHandler)) {
                throw new RuntimeException(sprintf('Cannot open file "%s" for writing.', $binary));
            }

            foreach ($this->httpClient->stream($response) as $chunk) {
                fwrite($fileHandler, $chunk->getContent());
            }

            fclose($fileHandler);

            $progressBar?->finish();
            $this->output->writeln('');

            chmod($binary, 0o777);
        }

        if (!is_file($binary)) {
            throw new RuntimeException(sprintf('Cannot find file "%s"', $binary));
        }

        return (new Process([$binary, ...$arguments], getcwd()))
            ->setTty(Process::isTtySupported())
            ->setEnv(['TERM' => 'xterm-256color', 'FORCE_COLOR' => '1'])
            ->setTimeout(null);
    }

    public function setOutput(
        ?SymfonyStyle $output,
    ): self {
        $this->output = $output;

        return $this;
    }
}
