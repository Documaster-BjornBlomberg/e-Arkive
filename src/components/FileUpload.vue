<script setup>
import { ref } from 'vue'

const file = ref(null)

const uploadFile = async () => {
  if (!file.value) return

  const formData = new FormData()
  formData.append('file', file.value)

  try {
    const response = await fetch('http://localhost:8080/upload', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) {
      throw new Error('Failed to upload file')
    }

    alert('File uploaded successfully!')
  } catch (error) {
    console.error(error)
    alert('Error uploading file')
  }
}
</script>

<template>
  <div>
    <h1>Upload a Document</h1>
    <input type="file" @change="(e) => (file.value = e.target.files[0])" />
    <button @click="uploadFile">Upload</button>
  </div>
</template>

<style scoped>
/* Add any necessary styles here */
</style>