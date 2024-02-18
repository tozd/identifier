import bs58 from "bs58"
import { parse } from "uuid"

const stringLength = 22
const bytesMinLength = 16
const idRegex = /^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{22}$/

export class Identifier {
  private readonly value: Uint8Array

  constructor(value: Uint8Array) {
    if (value.length !== 16) {
      throw new Error("invalid identifier")
    }
    this.value = value
  }

  // toString encodes Identifier value into a string using base 58 encoding.
  public toString(): string {
    const res = bs58.encode(this.value)
    if (res.length < stringLength) {
      // String might be shorter than stringLength to encode 128 bits, in that
      // we do zero left padding (character "1" in base 58).
      return "1".repeat(stringLength - res.length) + res
    }
    return res
  }

  public toJSON(): string {
    return this.toString()
  }

  // fromUUID returns the UUID encoded as an Identifier.
  public static fromUUID(data: string): Identifier {
    return Identifier.fromData(parse(data))
  }

  // fromData returns 16 bytes data encoded as an Identifier.
  public static fromData(data: Uint8Array): Identifier {
    return new Identifier(data)
  }

  // fromString parses a string-encoded identifier in base 58 encoding
  // into a corresponding Identifier value.
  public static fromString(data: string): Identifier {
    if (data.length !== stringLength) {
      throw new Error("invalid identifier")
    }
    const res = bs58.decode(data)
    if (res.length < bytesMinLength) {
      throw new Error("invalid identifier")
    }
    // String might longer than necessary to encode 128 bits, in that case we require extra bytes
    // at the beginning to be zero (or character "1" in base58), i.e., zero left padding.
    for (let i = 0; i + bytesMinLength < res.length; i++) {
      if (res[i] != 0) {
        throw new Error("invalid identifier")
      }
    }
    // We take the last 16 bytes.
    return new Identifier(res.slice(-16))
  }

  // new returns a new random identifier.
  public static new(): Identifier {
    const data = new Uint8Array(16)
    self.crypto.getRandomValues(data)
    return this.fromData(data)
  }
}

export function valid(id: string): boolean {
  if (!idRegex.test(id)) {
    return false
  }
  try {
    Identifier.fromString(id)
    return true
  } catch (error) {
    return false
  }
}
