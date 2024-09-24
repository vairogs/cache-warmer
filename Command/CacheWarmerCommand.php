<?php declare(strict_types = 1);

/*
 * This file is part of the Vairogs package.
 *
 * (c) Dāvis Zālītis (k0d3r1s) <davis@vairogs.com>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Vairogs\Component\CacheWarmer\Command;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\Process\Process;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;
use Vairogs\Component\CacheWarmer\CacheWarmerBinary;

use function array_pop;
use function explode;
use function gc_collect_cycles;
use function is_bool;
use function is_string;
use function sprintf;
use function trim;
use function usleep;

use const PHP_EOL;

#[AsCommand(
    name: 'vairogs:cache-warmer:watch',
)]
class CacheWarmerCommand extends Command
{
    private SymfonyStyle $io;

    public function __construct(
        private readonly CacheWarmerBinary $binary,
        ?string $name = null,
    ) {
        parent::__construct($name);
    }

    protected function configure(): void
    {
        $this
            ->addOption('cache', null, InputOption::VALUE_NONE, 'Clear cache instead of just warmup')
            ->addOption('exclude', null, InputOption::VALUE_OPTIONAL, 'Comma-separated directories not to watch', '')
            ->addOption('force', null, InputOption::VALUE_NONE, 'Force clear cache (rm -rf var/cache)')
            ->addOption('no-debug', null, InputOption::VALUE_NONE, 'Pass --no-debug to the symfony console')
            ->addOption('pools', null, InputOption::VALUE_OPTIONAL, 'Comma-separated list of pools to clear (can be true if passed without a value)', false)
            ->addOption('vendor', null, InputOption::VALUE_OPTIONAL, 'Comma-separated list of vendors to watch', '')
            ->addOption('manual', null, InputOption::VALUE_NONE, 'You have downloaded binary manually and want to skip download');
    }

    /**
     * @throws TransportExceptionInterface
     */
    protected function execute(
        InputInterface $input,
        OutputInterface $output,
    ): int {
        $this->binary->setOutput($this->io);

        $optionsToArguments = ['cache', 'force', 'no-debug', 'exclude', 'vendor', 'env', ];

        $arguments = [];

        foreach ($optionsToArguments as $argument) {
            $optionValue = $input->getOption($argument);

            if (is_bool($optionValue) && $optionValue) {
                $arguments[] = '--' . $argument;
            } elseif (is_string($optionValue) && '' !== trim($optionValue)) {
                $arguments[] = '--' . $argument . '=' . $optionValue;
            }
        }

        $pools = $input->getOption('pools');

        if (null === $pools) {
            $arguments[] = '--pools';
        } elseif (false !== $pools) {
            $arguments[] = '--pools=' . $pools;
        }

        $arguments[] = '.';

        ($process = $this->binary->createProcess($arguments, $input->getOption('manual')))->start();

        $outputBuffer = '';
        $counter = 0;

        while ($process->isRunning()) {
            $outputBuffer .= $process->getIncrementalOutput();
            $errorBuffer = $process->getIncrementalErrorOutput();

            $lines = explode(PHP_EOL, $outputBuffer);
            $outputBuffer = array_pop($lines);

            foreach ($lines as $line) {
                $this->io->writeln($line);
            }

            if (!empty($errorBuffer)) {
                foreach (explode(PHP_EOL, $errorBuffer) as $errorLine) {
                    if (!empty(trim($errorLine))) {
                        $this->io->error('Error: ' . $errorLine);
                    }
                }
            }

            if (600 /* 10 minutes */ === ++$counter) {
                gc_collect_cycles();
                $counter = 0;
            }

            usleep(1000000);
        }

        if (!empty($outputBuffer)) {
            $lines = explode(PHP_EOL, $outputBuffer);

            foreach ($lines as $line) {
                $this->io->writeln($line);
            }
        }

        if (!$process->isSuccessful()) {
            $this->io->error(sprintf('%s process has failed', CacheWarmerBinary::REPOSITORY));
            $this->cleanUp($process);

            return Command::FAILURE;
        }

        $this->io->success(sprintf('%s process has been stopped', CacheWarmerBinary::REPOSITORY));
        $this->cleanUp($process);

        return Command::SUCCESS;
    }

    protected function initialize(
        InputInterface $input,
        OutputInterface $output,
    ): void {
        $output->getFormatter()->setDecorated(true);
        $this->io = new SymfonyStyle($input, $output);
    }

    private function cleanUp(
        Process $process,
    ): void {
        if ($process->isRunning()) {
            $process->stop();
        }

        gc_collect_cycles();
    }
}
