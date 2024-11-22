export interface Shop {
  id: string
  name: string
  address: string
  type: string
  city: string
  location: {
    lat: number
    lng: number
  }
  openTime: string
  closeTime: string
}

export const shops: Shop[] = [
  {
    id: '1',
    name: '示例店铺',
    address: '北京市朝阳区xxx街xxx号',
    type: 'restaurant',
    city: 'beijing',
    location: {
      lat: 39.9042,
      lng: 116.4074
    },
    openTime: '08:00',
    closeTime: '22:00'
  }
  // 可以添加更多店铺数据
] 