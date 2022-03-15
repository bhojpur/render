import fs from 'fs';
import babel from '@rollup/plugin-babel';
import { terser } from 'rollup-plugin-terser';
import { sizeSnapshot } from 'rollup-plugin-size-snapshot';

import license from './utils/license-template.js';

const { version } = JSON.parse(fs.readFileSync(new URL('./package.json', import.meta.url)));

const input = './pkg/javascript/index.js';
const name = 'bhojpurRender';

const bannerPlugin = {
  banner: `/*!
@fileoverview bhojpur-render - A high performance rendering operations
@author Shashi Bhushan Rai
@version ${version}

${license}

*/`
}

export default [
  {
    input,
    output: { file: 'dist/bhojpur-render.js', format: 'umd', name },
    plugins: [
      bannerPlugin,
      babel()
    ]
  },
  {
    input,
    output: { file: 'dist/bhojpur-render-min.js', format: 'umd', name },
    plugins: [
      bannerPlugin,
      babel(),
      sizeSnapshot(),
      terser({
        output: { comments: /^!/ }
      })
    ]
  }
];