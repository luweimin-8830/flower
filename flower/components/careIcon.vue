<template>
  <image class="care-svg" :class="{ active }" :src="iconSrc" mode="aspectFit"
    :style="{ width: size + 'px', height: size + 'px' }" />
</template>

<script>
import { careIconSrc } from '../utils/careIcons'

export default {
  name: 'careIcon',
  props: {
    icon: { type: String, default: '' },
    type: { type: String, default: '' },
    name: { type: String, default: '' },
    active: { type: Boolean, default: false },
    size: { type: Number, default: 22 }
  },
  computed: {
    iconSrc() {
      return careIconSrc(this.icon, this.type, this.name)
    }
  }
}
</script>

<style scoped>
.care-svg { display: block; opacity: .78; }
.care-svg.active { opacity: 1; filter: brightness(0) invert(1); }
</style>
