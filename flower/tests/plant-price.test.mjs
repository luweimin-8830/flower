import { readFileSync } from 'node:fs';
import test from 'node:test';
import assert from 'node:assert/strict';
const source = readFileSync(new URL('../utils/plantPrice.js', import.meta.url));
const { parsePlantPrice, formatPlantPrice } = await import('data:text/javascript;base64,' + source.toString('base64'));
test('prices roundtrip exactly, with blank distinct from zero', () => {
 for (const [input, cents] of [['',null],['0',0],['0.01',1],['12.30',1230],['999999.99',99999999]]) {
  assert.equal(parsePlantPrice(input), cents);
  assert.equal(parsePlantPrice(formatPlantPrice(cents)), cents);
 }
 assert.equal(formatPlantPrice(undefined), '');
 for (const invalid of ['-1','1.234','1e3','1000000','abc']) assert.throws(() => parsePlantPrice(invalid));
});
