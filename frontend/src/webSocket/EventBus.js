import { ref } from 'vue';

const eventBus = ref({
  emit(event, data) {
    this[event] = data;
  },
});

export function useEventBus() {
  return eventBus;
}