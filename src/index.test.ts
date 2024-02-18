import { assert, test } from "vitest"

import { Identifier } from "./index"

// TODO: Convert to a fuzzing test and a benchmark.

test("Identifier.fromUUID", () => {
  for (let i = 0; i < 100000; i++) {
    const u = crypto.randomUUID()
    const i = Identifier.fromUUID(u)
    // @ts-expect-error: value is private
    assert.lengthOf(i.value, 16)
    const s = i.toString()
    assert.lengthOf(s, 22)
    assert.isTrue(Identifier.valid(s))
    const ii = Identifier.fromString(s)
    assert.deepEqual(i, ii)
  }
})

test("Identifier.new", () => {
  for (let i = 0; i < 100000; i++) {
    const i = Identifier.new()
    // @ts-expect-error: value is private
    assert.lengthOf(i.value, 16)
    const s = i.toString()
    assert.lengthOf(s, 22)
    assert.isTrue(Identifier.valid(s))
    const ii = Identifier.fromString(s)
    assert.deepEqual(i, ii)
  }
})

test.each([
  ["", false],
  ["42", false],
  ["CDEFGHJKLMNPQRSTUVWXYZ", true],
  ["zzzzzzzzzzzzzzzzzzzzzz", false],
  ["2222222222222222222222", true],
  ["2111111111111111111111", true],
  ["1111111111111111211111", true],
  ["1111111111111111111111", true],
])("Identifier.valid(%s)", (m, u) => {
  assert.equal(Identifier.valid(m), u)
})

test("Identifier.toJSON", () => {
  const i = Identifier.new()
  const o = { id: i }
  const data = JSON.stringify(o)
  assert.equal(data, `{"id":"${i.toString()}"}`)
  const obj = JSON.parse(data, (key, value) => {
    if (key === "id") {
      return Identifier.fromString(value)
    }
    return value
  })
  assert.deepEqual(obj, o)
})

test("Identifier.fromString error", () => {
  assert.throw(() => {
    Identifier.fromString("xxx")
  })
})
