export interface Shop {
  id: string
  type: string
  city: string
  location: {
    lat: number
    lng: number
  }
  openTime: string
  closeTime: string
}

// 生成随机位置
const generateRandomLocation = (baseLocation: { lat: number, lng: number }, radius: number) => {
  const randomAngle = Math.random() * Math.PI * 2
  const randomRadius = Math.random() * radius
  
  const lat = baseLocation.lat + (randomRadius * Math.cos(randomAngle) / 111)
  const lng = baseLocation.lng + (randomRadius * Math.sin(randomAngle) / (111 * Math.cos(baseLocation.lat * Math.PI / 180)))
  
  return { lat, lng }
}

// 城市数据
const cities = [
  { name: '北京', location: { lat: 39.9042, lng: 116.4074 } },
  { name: '上海', location: { lat: 31.2304, lng: 121.4737 } },
  { name: '广州', location: { lat: 23.1291, lng: 113.2644 } },
  { name: '深圳', location: { lat: 22.5431, lng: 114.0579 } },
  { name: '成都', location: { lat: 30.5728, lng: 104.0668 } }
]

// 店铺类型及其默认营业时间范围
const shopTypes = [
  { type: '餐厅', defaultOpen: '10:00', defaultClose: '22:00' },
  { type: '咖啡厅', defaultOpen: '08:00', defaultClose: '21:00' },
  { type: '美容院', defaultOpen: '09:00', defaultClose: '20:00' },
  { type: '健身房', defaultOpen: '06:00', defaultClose: '22:00' },
  { type: '超市', defaultOpen: '07:00', defaultClose: '22:00' },
  { type: '书店', defaultOpen: '09:00', defaultClose: '21:00' },
  { type: '服装店', defaultOpen: '10:00', defaultClose: '22:00' },
  { type: '电影院', defaultOpen: '10:00', defaultClose: '23:30' },
  { type: '药店', defaultOpen: '08:00', defaultClose: '21:00' },
  { type: '宠物店', defaultOpen: '09:00', defaultClose: '20:00' }
]

// 生成随机时间偏移
const generateTimeOffset = () => {
  const offsets = [-60, -30, 0, 30, 60] // 分钟偏移量
  return offsets[Math.floor(Math.random() * offsets.length)]
}

// 调整时间
const adjustTime = (time: string, offsetMinutes: number): string => {
  const [hours, minutes] = time.split(':').map(Number)
  const totalMinutes = hours * 60 + minutes + offsetMinutes
  const adjustedHours = Math.floor(totalMinutes / 60) % 24
  const adjustedMinutes = totalMinutes % 60
  return `${adjustedHours.toString().padStart(2, '0')}:${adjustedMinutes.toString().padStart(2, '0')}`
}

// 生成店铺数据
export const shops: Shop[] = []

cities.forEach(city => {
  // 每个城市生成20个店铺
  for (let i = 0; i < 20; i++) {
    const shopType = shopTypes[Math.floor(Math.random() * shopTypes.length)]
    const location = generateRandomLocation(city.location, 0.1) // 0.1度范围内随机分布
    
    // 生成略微不同的营业时间
    const openOffset = generateTimeOffset()
    const closeOffset = generateTimeOffset()
    const openTime = adjustTime(shopType.defaultOpen, openOffset)
    const closeTime = adjustTime(shopType.defaultClose, closeOffset)

    shops.push({
      id: `${city.name}-${shopType.type}-${i}`,
      type: shopType.type,
      city: city.name,
      location,
      openTime,
      closeTime
    })
  }
})

// 确保数据已经生成
console.log(`Generated ${shops.length} shops`) 