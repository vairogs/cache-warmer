<?php declare(strict_types = 1);

/*
 * This file is part of the Vairogs package.
 *
 * (c) Dāvis Zālītis (k0d3r1s) <davis@vairogs.com>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Vairogs\Component\CacheWarmer\DependencyInjection;

use Symfony\Component\Config\Definition\Builder\ArrayNodeDefinition;
use Vairogs\Bundle\DependencyInjection\AbstractDependencyConfiguration;
use Vairogs\Component\Functions\Vairogs;

use function sprintf;

class CacheWarmerConfiguration extends AbstractDependencyConfiguration
{
    public function addSection(
        ArrayNodeDefinition $rootNode,
        callable $enableIfStandalone,
        string $component,
    ): void {
        $rootNode
            ->children()
                ->arrayNode($component)
                    ->{$enableIfStandalone(sprintf('%s/%s', Vairogs::VAIROGS, $component), self::class)}()
                    ->children()
                        ->scalarNode('version')->defaultValue('latest')->end()
                    ->end()
                ->end();
    }
}
