<script setup>
import { ref } from 'vue'

const file = ref(null)
const metadataList = ref([{ key: '', value: '' }])

const addMetadata = () => {
  metadataList.value.push({ key: '', value: '' })
}

const removeMetadata = (index) => {
  metadataList.value.splice(index, 1)
}

const uploadFile = async () => {
  if (!file.value) {
    alert('Please select a file to upload')
    return
  }

  const reader = new FileReader()
  reader.onload = async (e) => {
    const base64Data = e.target.result.split(',')[1]

    const metadata = metadataList.value.filter(meta => meta.key && meta.value)

    const fileInput = {
      name: file.value.name,
      size: file.value.size,
      contentType: file.value.type,
      fileData: base64Data,
      metadata,
    }

    const query = `
      mutation($input: FileInput!) {
        saveFile(input: $input) {
          id
          name
        }
      }
    `

    const variables = { input: fileInput }

    try {
      const response = await fetch('http://localhost:8080/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query, variables }),
      })

      const result = await response.json()
      if (result.errors) {
        throw new Error(result.errors[0].message)
      }

      alert('File uploaded successfully!')
    } catch (error) {
      alert('Error uploading file: ' + error.message)
    }
  }

  reader.readAsDataURL(file.value)
}

const onFileChange = (e) => {
  file.value = e.target.files[0]
}
</script>

<template>
  <div>
    <h1>Upload a Document</h1>
    <input type="file" @change="onFileChange" />

    <div v-for="(meta, index) in metadataList" :key="index">
      <label>Metadata Key:</label>
      <input type="text" v-model="meta.key" />

      <label>Metadata Value:</label>
      <input type="text" v-model="meta.value" />

      <button @click="removeMetadata(index)">Remove</button>
    </div>

    <button @click="addMetadata">Add Metadata</button>
    <button @click="uploadFile">Upload</button>
  </div>
</template>

<style scoped>
/* Add styles if needed */
</style>