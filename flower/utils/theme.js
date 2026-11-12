export const THEME_STORAGE_KEY = 'gardenTheme'

export const themes = {
  mint: {
    name: '薄荷晨露',
    swatch: '#72A990',
    vars: {
      '--page-bg': '#F1F7F4',
      '--card-bg': '#FCFFFD',
      '--surface-bg': '#E8F1ED',
      '--border-color': '#D6E5DE',
      '--primary-color': '#72A990',
      '--primary-strong': '#4F8B70',
      '--primary-soft': '#DCEDE4',
      '--text-color': '#203029',
      '--text-sub': '#6C7A73',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(49, 96, 73, 0.10)'
    }
  },
  apricot: {
    name: '杏陶暖阳',
    swatch: '#C6815E',
    vars: {
      '--page-bg': '#FBF4EC',
      '--card-bg': '#FFFDF9',
      '--surface-bg': '#F4E9DF',
      '--border-color': '#EADCCF',
      '--primary-color': '#C6815E',
      '--primary-strong': '#A96345',
      '--primary-soft': '#F3DED1',
      '--text-color': '#342620',
      '--text-sub': '#7D6D64',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(135, 81, 52, 0.11)'
    }
  },
  lavender: {
    name: '紫藤轻雾',
    swatch: '#8D7EAA',
    vars: {
      '--page-bg': '#F6F4F9',
      '--card-bg': '#FEFDFF',
      '--surface-bg': '#EEEAF3',
      '--border-color': '#E1DCE9',
      '--primary-color': '#8D7EAA',
      '--primary-strong': '#70618F',
      '--primary-soft': '#E7E0F0',
      '--text-color': '#292632',
      '--text-sub': '#746E7D',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(85, 70, 116, 0.11)'
    }
  },
  lightBlue: {
    name: '晴空微澜',
    swatch: '#72A7C5',
    vars: {
      '--page-bg': '#F1F7FB',
      '--card-bg': '#FCFEFF',
      '--surface-bg': '#E7F0F7',
      '--border-color': '#D6E4EE',
      '--primary-color': '#72A7C5',
      '--primary-strong': '#4F83A3',
      '--primary-soft': '#DCECF5',
      '--text-color': '#20303A',
      '--text-sub': '#6C7A84',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(49, 86, 110, 0.10)'
    }
  },
  flamingoRed: {
    name: '火烈鸟红',
    swatch: '#EF8E87',
    vars: {
      '--page-bg': '#FFF3F2',
      '--card-bg': '#FFFCFC',
      '--surface-bg': '#FCE5E3',
      '--border-color': '#F1CECB',
      '--primary-color': '#EF8E87',
      '--primary-strong': '#C85E59',
      '--primary-soft': '#F8D7D4',
      '--text-color': '#3C2928',
      '--text-sub': '#846B69',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(199, 91, 84, 0.13)'
    }
  },
  forestNight: {
    name: '松林夜语',
    swatch: '#4E8069',
    dark: true,
    vars: {
      '--page-bg': '#101B16',
      '--card-bg': '#192720',
      '--surface-bg': '#14221B',
      '--border-color': '#2B4035',
      '--primary-color': '#69A88A',
      '--primary-strong': '#91C7AD',
      '--primary-soft': '#29483A',
      '--text-color': '#ECF4F0',
      '--text-sub': '#A6B9AF',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(0, 0, 0, 0.30)'
    }
  },
  deepOcean: {
    name: '深海星潮',
    swatch: '#4F86B4',
    dark: true,
    vars: {
      '--page-bg': '#0E1823',
      '--card-bg': '#182738',
      '--surface-bg': '#132130',
      '--border-color': '#2A4055',
      '--primary-color': '#629BC8',
      '--primary-strong': '#93C1E3',
      '--primary-soft': '#274762',
      '--text-color': '#EDF5FA',
      '--text-sub': '#A2B5C4',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(0, 0, 0, 0.32)'
    }
  },
  violetNight: {
    name: '暮紫微光',
    swatch: '#806EAA',
    dark: true,
    vars: {
      '--page-bg': '#191621',
      '--card-bg': '#272333',
      '--surface-bg': '#211D2C',
      '--border-color': '#40394F',
      '--primary-color': '#9582BD',
      '--primary-strong': '#C0B0DF',
      '--primary-soft': '#433858',
      '--text-color': '#F3EFF8',
      '--text-sub': '#B7ADBF',
      '--button-text': '#FFFFFF',
      '--shadow-color': 'rgba(0, 0, 0, 0.32)'
    }
  }
}

export function normalizeTheme(key) {
  return themes[key] ? key : 'mint'
}

export function getStoredTheme() {
  return normalizeTheme(uni.getStorageSync(THEME_STORAGE_KEY))
}

export function storeTheme(key) {
  const next = normalizeTheme(key)
  uni.setStorageSync(THEME_STORAGE_KEY, next)
  return next
}

export function getThemeStyle(key) {
  return themes[normalizeTheme(key)].vars
}

export const themeMixin = {
  data() {
    return { themeKey: getStoredTheme() }
  },
  computed: {
    themeStyle() {
      return getThemeStyle(this.themeKey)
    },
    activeTheme() {
      return themes[this.themeKey] || themes.mint
    }
  },
  created() {
    this.__themeChangeHandler = (key) => {
      this.themeKey = normalizeTheme(key)
    }
    uni.$on('garden-theme-change', this.__themeChangeHandler)
  },
  beforeUnmount() {
    if (this.__themeChangeHandler) uni.$off('garden-theme-change', this.__themeChangeHandler)
  },
  methods: {
    syncGardenTheme() {
      const stored = getStoredTheme()
      if (this.themeKey !== stored) this.themeKey = stored
      return stored
    },
    setGardenTheme(key) {
      const next = storeTheme(key)
      this.themeKey = next
      uni.$emit('garden-theme-change', next)
      uni.setNavigationBarColor({
        frontColor: themes[next].dark ? '#ffffff' : '#000000',
        backgroundColor: themes[next].vars['--page-bg'],
        animation: { duration: 200, timingFunc: 'easeIn' }
      })
    }
  }
}
