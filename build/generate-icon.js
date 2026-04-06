const sharp = require('sharp')
const path = require('path')
const fs = require('fs')

const svgPath = path.join(__dirname, 'appicon.svg')
const svg = fs.readFileSync(svgPath)

async function generate() {
  // PNG 512x512 (для Wails)
  await sharp(svg)
    .resize(512, 512)
    .png()
    .toFile(path.join(__dirname, 'appicon.png'))
  console.log('✓ appicon.png')

  // Несколько размеров для ICO
  const sizes = [16, 32, 48, 64, 128, 256]
  const pngBuffers = await Promise.all(
    sizes.map(s => sharp(svg).resize(s, s).png().toBuffer())
  )

  // Генерируем ICO вручную (Multi-size ICO format)
  const ico = buildIco(pngBuffers, sizes)
  fs.writeFileSync(path.join(__dirname, 'windows', 'icon.ico'), ico)
  console.log('✓ windows/icon.ico')
}

function buildIco(pngBuffers, sizes) {
  const n = sizes.length
  const headerSize = 6
  const dirEntrySize = 16
  const dataOffset = headerSize + dirEntrySize * n

  let totalDataSize = 0
  const offsets = []
  for (const buf of pngBuffers) {
    offsets.push(dataOffset + totalDataSize)
    totalDataSize += buf.length
  }

  const ico = Buffer.alloc(dataOffset + totalDataSize)

  // ICONDIR header
  ico.writeUInt16LE(0, 0)       // Reserved
  ico.writeUInt16LE(1, 2)       // Type: ICO
  ico.writeUInt16LE(n, 4)       // Count

  // ICONDIRENTRY for each image
  for (let i = 0; i < n; i++) {
    const base = headerSize + dirEntrySize * i
    const s = sizes[i]
    ico.writeUInt8(s >= 256 ? 0 : s, base)      // Width (0 = 256)
    ico.writeUInt8(s >= 256 ? 0 : s, base + 1)  // Height
    ico.writeUInt8(0, base + 2)   // Color count
    ico.writeUInt8(0, base + 3)   // Reserved
    ico.writeUInt16LE(1, base + 4) // Planes
    ico.writeUInt16LE(32, base + 6) // Bit count
    ico.writeUInt32LE(pngBuffers[i].length, base + 8)  // Size
    ico.writeUInt32LE(offsets[i], base + 12)            // Offset
  }

  // PNG data
  let pos = dataOffset
  for (const buf of pngBuffers) {
    buf.copy(ico, pos)
    pos += buf.length
  }

  return ico
}

generate().catch(console.error)
