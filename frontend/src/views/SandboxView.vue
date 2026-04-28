<template>
  <v-container class="py-8">
    <v-row>
      <!-- Controls -->
      <v-col cols="12" md="4" lg="3">
        <v-card class="pa-6 mb-6 overflow-hidden">
          <h3 class="text-h5 font-weight-bold mb-6 d-flex align-center">
            <v-icon color="primary" class="mr-2">mdi-cog-outline</v-icon>
            Параметры
          </h3>

          <v-select
            v-model="method"
            :items="methods"
            label="Метод интерполяции"
            variant="outlined"
            density="comfortable"
            item-title="text"
            item-value="value"
            class="mb-4"
          ></v-select>

          <v-divider class="my-6"></v-divider>

          <div class="d-flex justify-space-between align-center mb-4">
            <span class="text-subtitle-2 font-weight-bold text-text-secondary">Узлы (X, Y)</span>
            <v-btn
              icon="mdi-plus"
              size="x-small"
              color="primary"
              variant="flat"
              @click="addPoint"
            ></v-btn>
          </div>

          <v-scroll-y-transition group>
            <div v-for="(p, index) in points" :key="index" class="mb-3">
              <v-row dense align="center">
                <v-col cols="5">
                  <v-text-field
                    v-model.number="p.x"
                    type="number"
                    label="X"
                    variant="outlined"
                    density="compact"
                    hide-details
                    class="math-font"
                  ></v-text-field>
                </v-col>
                <v-col cols="5">
                  <v-text-field
                    v-model.number="p.y"
                    type="number"
                    label="Y"
                    variant="outlined"
                    density="compact"
                    hide-details
                    class="math-font"
                  ></v-text-field>
                </v-col>
                <v-col cols="2" class="text-right">
                  <v-btn
                    icon="mdi-close-circle-outline"
                    variant="text"
                    color="error"
                    size="small"
                    @click="removePoint(index)"
                    :disabled="points.length <= 2"
                  ></v-btn>
                </v-col>
              </v-row>
            </div>
          </v-scroll-y-transition>

          <v-divider class="my-6"></v-divider>

          <v-text-field
            v-model.number="targetX"
            type="number"
            label="Целевое значение X"
            variant="outlined"
            density="comfortable"
            class="math-font mb-6"
            persistent-hint
            hint="Точка, в которой нужно вычислить значение"
          ></v-text-field>

          <v-card v-if="result !== null" color="primary" variant="tonal" class="rounded-lg border-0">
            <v-card-text class="text-center py-4">
              <div class="text-caption text-uppercase font-weight-bold mb-1 opacity-70">Результат f(x)</div>
              <div class="text-h3 font-weight-black math-font">{{ result.toFixed(4) }}</div>
              <v-chip v-if="cached" size="x-small" color="primary" variant="flat" class="mt-2">ИЗ КЭША</v-chip>
            </v-card-text>
          </v-card>
        </v-card>
      </v-col>

      <!-- Chart -->
      <v-col cols="12" md="8" lg="9">
        <v-card class="pa-6 h-100 d-flex flex-column">
          <h3 class="text-h5 font-weight-bold mb-6 d-flex align-center">
            <v-icon color="secondary" class="mr-2">mdi-chart-bell-curve</v-icon>
            Визуализация
          </h3>
          
          <v-alert v-if="error" type="error" variant="tonal" class="mb-4 rounded-lg">
            {{ error }}
          </v-alert>

          <div class="flex-grow-1 relative-container">
            <div v-if="loading" class="overlay-loader d-flex justify-center align-center">
              <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
            </div>
            <ChartView
              :points="points"
              :curve="curve"
              :targetPoint="result !== null ? { x: targetX, y: result } : null"
            />
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- History -->
    <v-row v-if="history && history.length > 0">
      <v-col cols="12">
        <v-card class="mt-8 pa-6">
          <div class="d-flex justify-space-between align-center mb-6">
            <h3 class="text-h5 font-weight-bold d-flex align-center">
              <v-icon color="text-secondary" class="mr-2">mdi-history</v-icon>
              История вычислений
            </h3>
            <v-btn
              prepend-icon="mdi-delete-sweep-outline"
              variant="text"
              color="error"
              class="text-none font-weight-bold"
              @click="clearHistory"
            >Очистить историю</v-btn>
          </div>

          <v-table hover class="history-table">
            <thead>
              <tr>
                <th class="text-uppercase text-caption font-weight-bold">Дата</th>
                <th class="text-uppercase text-caption font-weight-bold">Метод</th>
                <th class="text-uppercase text-caption font-weight-bold">Цель X</th>
                <th class="text-uppercase text-caption font-weight-bold">Результат Y</th>
                <th class="text-right text-uppercase text-caption font-weight-bold">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="h in history" :key="h.id">
                <td class="text-body-2 text-text-secondary">{{ new Date(h.created_at).toLocaleString() }}</td>
                <td>
                  <v-chip :color="getMethodColor(h.method)" size="x-small" variant="flat" class="font-weight-bold">
                    {{ getMethodText(h.method) }}
                  </v-chip>
                </td>
                <td class="math-font">{{ h.target_x }}</td>
                <td class="math-font font-weight-bold">{{ h.result.toFixed(4) }}</td>
                <td class="text-right">
                  <v-btn
                    icon="mdi-restore"
                    variant="text"
                    color="secondary"
                    size="small"
                    @click="restore(h)"
                    title="Восстановить параметры"
                  ></v-btn>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import axios from 'axios'
import debounce from 'lodash/debounce'
import ChartView from '../components/ChartView.vue'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const method = ref('linear')
const methods = [
  { text: 'Линейная', value: 'linear' },
  { text: 'Лагранж', value: 'lagrange' },
  { text: 'Ньютон', value: 'newton' }
]

const points = ref([
  { x: 0, y: 0 },
  { x: 1, y: 1 },
  { x: 2, y: 4 }
])

const targetX = ref(1.5)
const result = ref(null)
const curve = ref([])
const cached = ref(false)
const loading = ref(false)
const error = ref(null)
const history = ref([])

const addPoint = () => {
  const lastX = points.value[points.value.length - 1].x
  points.value.push({ x: Number(lastX) + 1, y: 0 })
}

const removePoint = (index) => {
  points.value.splice(index, 1)
}

const calculate = async () => {
  if (!points.value || points.value.length < 2) return

  loading.value = true
  error.value = null
  try {
    const formattedPoints = points.value.map(p => ({ 
      x: parseFloat(p.x) || 0, 
      y: parseFloat(p.y) || 0 
    }))
    
    const resp = await axios.post(`${API_URL}/interpolate`, {
      method: method.value,
      points: formattedPoints,
      target_x: parseFloat(targetX.value) || 0
    })
    
    if (resp.data) {
      result.value = resp.data.result
      // Убедимся, что кривая - это массив
      curve.value = Array.isArray(resp.data.curve) ? resp.data.curve : []
      cached.value = resp.data.cached
      fetchHistory()
    }
  } catch (err) {
    console.error('Calculation error:', err)
    error.value = err.response?.data?.error || 'Ошибка при вычислении'
    result.value = null
    curve.value = []
  } finally {
    loading.value = false
  }
}

const debouncedCalculate = debounce(calculate, 300)

const fetchHistory = async () => {
  try {
    const resp = await axios.get(`${API_URL}/history`)
    history.value = Array.isArray(resp.data) ? resp.data : []
  } catch (err) {
    console.error('Failed to fetch history', err)
    history.value = []
  }
}

const clearHistory = async () => {
  if (!confirm('Вы уверены, что хотите очистить всю историю?')) return
  try {
    await axios.delete(`${API_URL}/history`)
    history.value = []
  } catch (err) {
    console.error('Failed to clear history', err)
  }
}

const restore = (h) => {
  method.value = h.method
  points.value = JSON.parse(JSON.stringify(h.points))
  targetX.value = h.target_x
  calculate()
}

const getMethodText = (m) => methods.find(opt => opt.value === m)?.text || m
const getMethodColor = (m) => {
  if (m === 'linear') return 'primary'
  if (m === 'lagrange') return 'secondary'
  return 'warning'
}

watch([method, points, targetX], () => {
  debouncedCalculate()
}, { deep: true })

onMounted(() => {
  calculate()
  fetchHistory()
})
</script>

<style scoped>
.relative-container {
  position: relative;
  min-height: 400px;
}
.overlay-loader {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.7);
  z-index: 10;
  backdrop-filter: blur(2px);
}
.history-table th {
  background-color: #fcfcfd !important;
}
</style>
