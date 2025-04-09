<template>
  <div class="main-layout">
    <HeaderSection />
    <main>
      <slot />
    </main>
  </div>
</template>

<script setup>
import HeaderSection from '@/components/organisms/HeaderSection.vue'
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

main {
  flex: 1;
  padding: 20px;
}
</style>