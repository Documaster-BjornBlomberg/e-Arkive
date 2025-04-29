import { ref } from 'vue';

export function useFileHandling() {
  const file = ref<File | null>(null);

  const onFileChange = (e: Event): void => {
    const target = e.target as HTMLInputElement;
    if (target && target.files && target.files.length > 0) {
      file.value = target.files[0];
    }
  };

  return {
    file,
    onFileChange,
  };
}
