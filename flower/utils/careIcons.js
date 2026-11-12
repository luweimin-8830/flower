export const CARE_ICON_OPTIONS = [
  { key: 'care-water', label: '浇水' },
  { key: 'care-fertilize', label: '施肥' },
  { key: 'care-repot', label: '换土' },
  { key: 'care-sprout', label: '发芽' },
  { key: 'care-spray', label: '喷药' },
  { key: 'care-prune', label: '修剪' },
  { key: 'care-pest', label: '除虫' },
  { key: 'care-dust', label: '除尘' },
  { key: 'care-loosen', label: '松土' },
  { key: 'care-sun', label: '晒太阳' },
  { key: 'care-record', label: '记录' },
  { key: 'care-ventilate', label: '通风' }
]

const LEGACY_ICON_MAP = {
  'plant-jiaoshui1': 'care-water',
  'plant-a-jiaoshui23': 'care-water',
  'plant-shifei1': 'care-fertilize',
  'plant-dianpu': 'care-repot',
  'plant-kongzhonghuayuan': 'care-sprout',
  'plant-a-xiujian13': 'care-spray',
  'plant-a-xiujian23': 'care-prune',
  'plant-chuchong1': 'care-pest',
  'plant-chuchong2': 'care-pest',
  'plant-shequhuodong': 'care-sun',
  'plant-a-Frame9': 'care-record',
  'plant-a-Frame1': 'care-ventilate',
  'plant-zhiwuzhiyuan-duorouzhiwuyuan': 'care-record'
}

const KEYWORDS = [
  ['care-water', /浇水|补水|water/i],
  ['care-fertilize', /施肥|肥料|fertiliz/i],
  ['care-repot', /换土|换盆|土壤|soil|repot/i],
  ['care-sprout', /发芽|生根|育苗|sprout|seed/i],
  ['care-spray', /喷药|喷雾|药剂|spray/i],
  ['care-prune', /修剪|剪枝|prun|trim/i],
  ['care-pest', /除虫|杀虫|虫害|pest|bug/i],
  ['care-dust', /除尘|擦叶|清洁|dust|clean/i],
  ['care-loosen', /松土|翻土|loosen|rake/i],
  ['care-sun', /晒太阳|日照|补光|sun|light/i],
  ['care-ventilate', /通风|透气|ventilat|air/i],
  ['care-record', /记录|成长|观察|record|growth|note/i]
]

export function normalizeCareIcon(icon, type, name) {
  if (CARE_ICON_OPTIONS.some(item => item.key === icon)) return icon
  if (LEGACY_ICON_MAP[icon]) return LEGACY_ICON_MAP[icon]
  const text = `${type || ''} ${name || ''}`
  const matched = KEYWORDS.find(([, pattern]) => pattern.test(text))
  return matched ? matched[0] : 'care-record'
}

export function careIconSrc(icon, type, name) {
  const key = normalizeCareIcon(icon, type, name).replace('care-', '')
  return `/static/care-icons/${key}.svg`
}
