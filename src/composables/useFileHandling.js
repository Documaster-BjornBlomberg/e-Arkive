import { ref } from 'vue';

export function useFileHandling() {
  const file = ref(null);

  const onFileChange = (e) => {
    file.value = e.target.files[0];
  };

  return {
    file,
    onFileChange,
  };
}
