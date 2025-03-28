<script setup>
import { ref, onMounted } from 'vue'

const file = ref(null)
const metadata = ref({ key: '', value: '' })
const files = ref([])

const log = (message) => console.log(`[FileUpload] ${message}`)

const fetchFiles = async () => {
  log("Fetching files from the server...")
  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query: '{ getFiles { id name size contentType createdAt metadata { key value } } }' }),
    })

    if (!response.ok) {
      throw new Error('Failed to fetch files')
    }

    const result = await response.json()
    files.value = result.data.getFiles
    log("Files fetched successfully")
  } catch (error) {
    log(`Error fetching files: ${error.message}`)
    alert('Error fetching files')
  }
}

const uploadFile = async () => {
  if (!file.value) {
    log("No file selected for upload")
    return
  }

  log(`Uploading file: ${file.value.name}`)

  const query = `
    mutation($input: FileInput!) {
      saveFile(input: $input) {
        id
        name
      }
    }
  `

  const variables = {
    input: {
      name: file.value.name,
      size: file.value.size,
      contentType: file.value.type,
      metadata: [{ key: metadata.value.key, value: metadata.value.value }],
    },
  }

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables }),
    })

    if (!response.ok) {
      throw new Error('Failed to upload file')
    }

    log("File uploaded successfully")
    alert('File uploaded successfully!')
    fetchFiles()
  } catch (error) {
    log(`Error uploading file: ${error.message}`)
    alert('Error uploading file')
  }
}

onMounted(() => {
  fetchFiles()
})
</script>

<template>
  <div>
    <h1>Upload a Document</h1>
    <input type="file" @change="(e) => (file.value = e.target.files[0])" />
    <div>
      <label>Metadata Key:</label>
      <input type="text" v-model="metadata.key" />
    </div>
    <div>
      <label>Metadata Value:</label>
      <input type="text" v-model="metadata.value" />
    </div>
    <button @click="uploadFile">Upload</button>

    <h2>Uploaded Files</h2>
    <ul>
      <li v-for="file in files" :key="file.id">
        <strong>{{ file.name }}</strong> ({{ file.size }} bytes, {{ file.contentType }})
        <ul>
          <li v-for="meta in file.metadata" :key="meta.key">
            {{ meta.key }}: {{ meta.value }}
          </li>
        </ul>
      </li>
    </ul>
  </div>
</template>

<style scoped>
/* Add any necessary styles here */
</style>