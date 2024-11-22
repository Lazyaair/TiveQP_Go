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

// 模拟数据
export const shops: Shop[] = [
  {
    id: '1',
    type: 'Hair Salons',
    city: 'ORLANDO',
    location: {
      lat: 28.4498736,
      lng: -81.4863524
    },
    openTime: '10:00',
    closeTime: '18:00'
  }
] 