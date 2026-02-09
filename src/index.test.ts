import { assert, test } from "vitest"

import { Identifier } from "./index.ts"

// TODO: Convert to a fuzzing test and a benchmark.

test("Identifier.fromUUID", () => {
  for (let i = 0; i < 10000; i++) {
    const u = crypto.randomUUID()
    const id = Identifier.fromUUID(u)
    // @ts-expect-error: value is private
    assert.lengthOf(id.value, 16)
    const s = id.toString()
    assert.lengthOf(s, 22)
    assert.isTrue(Identifier.valid(s))
    const ii = Identifier.fromString(s)
    assert.deepEqual(id, ii)
  }
})

test("Identifier.new", () => {
  for (let i = 0; i < 10000; i++) {
    const id = Identifier.new()
    // @ts-expect-error: value is private
    assert.lengthOf(id.value, 16)
    const s = id.toString()
    assert.lengthOf(s, 22)
    assert.isTrue(Identifier.valid(s))
    const ii = Identifier.fromString(s)
    assert.deepEqual(id, ii)
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
  const obj = JSON.parse(data, (key: string, value: unknown) => {
    if (key === "id") {
      return Identifier.fromString(value as string)
    }
    return value
  }) as unknown
  assert.deepEqual(obj, o)
})

test("Identifier.fromString error", () => {
  assert.throw(() => {
    Identifier.fromString("xxx")
  })
})

test("Identifier.from - produces valid identifier", async () => {
  const id = await Identifier.from("test")
  // @ts-expect-error: value is private
  assert.lengthOf(id.value, 16)
  const s = id.toString()
  assert.lengthOf(s, 22)
  assert.isTrue(Identifier.valid(s))
  assert.equal(s, "LhYXZThRsu1RXG5ddRZiUt")
})

test("Identifier.from - deterministic", async () => {
  const i1 = await Identifier.from("value1", "value2", "value3")
  const i2 = await Identifier.from("value1", "value2", "value3")
  assert.deepEqual(i1, i2)
  assert.equal(i1.toString(), "UhsHqGT45sDirscEZLxmC3")
})

test("Identifier.from - different inputs produce different outputs", async () => {
  const i1 = await Identifier.from("value1")
  const i2 = await Identifier.from("value2")
  assert.notDeepEqual(i1, i2)
  assert.equal(i1.toString(), "8UwJ16f3LZEDo1EWEPR1Ua")
  assert.equal(i2.toString(), "1eNbijZLjE6RCP9J3v6yz1")
})

test("Identifier.from - order matters", async () => {
  const i1 = await Identifier.from("value1", "value2")
  const i2 = await Identifier.from("value2", "value1")
  assert.notDeepEqual(i1, i2)
  assert.equal(i1.toString(), "ReKxivb3BXqpCurBhx657A")
  assert.equal(i2.toString(), "FtRyFysugNvJjVkMztmAuW")
})

test("Identifier.from - single vs multiple values", async () => {
  const i1 = await Identifier.from("value1value2")
  const i2 = await Identifier.from("value1", "value2")
  assert.notDeepEqual(i1, i2)
  assert.equal(i1.toString(), "21DW9wo4kBwGXVxbPW69oQ")
  assert.equal(i2.toString(), "ReKxivb3BXqpCurBhx657A")
})

test("Identifier.from - empty string", async () => {
  const id = await Identifier.from("")
  // @ts-expect-error: value is private
  assert.lengthOf(id.value, 16)
  const s = id.toString()
  assert.isTrue(Identifier.valid(s))
  assert.equal(s, "V7jseQevszwMPhi4evidTR")
})

test("Identifier.from - multiple empty strings", async () => {
  const i1 = await Identifier.from("", "")
  const i2 = await Identifier.from("")
  assert.notDeepEqual(i1, i2)
  assert.equal(i1.toString(), "Cbyu7w2KmnA6ZJVbgsHpHH")
  assert.equal(i2.toString(), "V7jseQevszwMPhi4evidTR")
})

test("Identifier.from - unicode normalization NFC", async () => {
  // U+00E9 (é) vs U+0065 U+0301 (e + combining acute accent)
  // Both should normalize to U+00E9 in NFC.
  const i1 = await Identifier.from("\u00e9") // é as single character.
  const i2 = await Identifier.from("\u0065\u0301") // e + combining acute.
  assert.deepEqual(i1, i2, "NFC normalization should make these equal")
  assert.equal(i1.toString(), "ADHVGvUx5PLGDsjwjY5BA9")
})

test("Identifier.from - specific known values", async () => {
  // Test with known values to ensure consistency across implementations.
  const id = await Identifier.from("test", "value")
  const s = id.toString()
  // This is just to verify the implementation produces consistent results.
  assert.lengthOf(s, 22)
  assert.isTrue(Identifier.valid(s))
  assert.equal(s, "J1oVAcLajL9m5GgBJ1eeqz")

  // Same input should always produce same output.
  const i2 = await Identifier.from("test", "value")
  assert.deepEqual(id, i2)
})

test("Identifier.from - cascading hash", async () => {
  // Each additional value should change the hash.
  const i1 = await Identifier.from("a")
  const i2 = await Identifier.from("a", "b")
  const i3 = await Identifier.from("a", "b", "c")
  assert.notDeepEqual(i1, i2)
  assert.notDeepEqual(i2, i3)
  assert.notDeepEqual(i1, i3)
  assert.equal(i1.toString(), "S1yrYnjHbfbiTySsN9h1eC")
  assert.equal(i2.toString(), "KorZ8VDpKQvHrZd2njXraU")
  assert.equal(i3.toString(), "469q6wDNXV222gSefVXCrJ")
})
