import { readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import type { Plugin } from 'vite';

const __filename = fileURLToPath(import.meta.url);
// 创建 esm 的dirname
const __dirname = dirname(__filename);

export function lumaBoot(): Plugin {
  return {
    name: 'luma-boot',
    transformIndexHtml: {
      order: 'pre',
      handler() {
        return [
          {
            tag: 'script',
            children: readFileSync(
              resolve(__dirname, '../src/boot/theme-init.js'),
              'utf-8',
            ),
          },
          {
            tag: 'style',
            children: readFileSync(
              resolve(__dirname, '../src/boot/luma-boot.css'),
              'utf-8',
            ),
          },
        ];
      },
    },
  };
}
