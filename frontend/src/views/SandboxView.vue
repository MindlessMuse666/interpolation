<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="4">
        <v-card class="pa-4 mb-4">
          <v-card-title>Параметры</v-card-title>
          <v-select
            v-model="method"
            :items="methods"
            label="Метод интерполяции"
            item-title="text"
            item-value="value"
          ></v-select>

          <v-divider class="my-4"></v-divider>

          <v-list density="compact">
            <v-list-subheader>Точки (X, Y)</v-list-subheader>
            <v-list-item v-for="(p, index) in points" :key="index">
              <v-row dense>
                <v-col cols="5">
                  <v-text-field
                    v-model.number="p.x"
                    type="number"
                    label="X"
                    hide-details
                    density="compact"
                  ></v-text-field>
                </v-col>
                <v-col cols="5">
                  <v-text-field
                    v-model.number="p.y"
                    type="number"
                    label="Y"
                    hide-details
                    density="compact"
                  ></v-text-field>
                </v-col>
                <v-col cols="2">
                  <v-btn
                    icon="mdi-delete"
                    variant="text"
                    color="error"
                    @click="removePoint(index)"
                    :disabled="points.length <= 2"
                  ></v-btn>
                </v-col>
              </v-row>
            </v-list-item>
          </v-list>
          <v-btn
            block
            prepend-icon="mdi-plus"
            variant="tonal"
            class="mt-2"
            @click="addPoint"
          >Добавить точку</v-btn>

          <v-divider class="my-4"></v-divider>

          <v-text-field
            v-model.number="targetX"
            type="number"
            label="Целевое значение X"
          ></v-text-field>
        </v-card>

        <v-card v-if="result !== null" color="success" theme="dark">
          <v-card-text class="text-center">
            <div class="text-h6">Результат:</div>
            <div class="text-h4">{{ result.toFixed(4) }}</div>
            <div v-if="cached" class="text-caption mt-2">Из кэша</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="pa-4 h-100">
          <v-card-title>График</v-card-title>
          <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
          <div v-if="loading" class="d-flex justify-center align-center h-100">
            <v-progress-circular indeterminate></v-progress-circular>
          </div>
          <ChartView
            v-else
            :points="points"
            :curve="curve"
            :targetPoint="result !== null ? { x: targetX, y: result } : null"
          />
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card class="mt-4">
          <v-card-title>История вычислений</v-card-title>
          <v-table>
            <thead>
              <tr>
                <th>Дата</th>
                <th>Метод</th>
                <th>X</th>
                <th>Результат</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="h in history" :key="h.id">
                <td>{{ new Date(h.created_at).toLocaleString() }}</td>
                <td>{{ h.method }}</td>
                <td>{{ h.target_x }}</td>
                <td>{{ h.result.toFixed(4) }}</td>
                <td>
                  <v-btn
                    size="small"
                    variant="text"
                    color="primary"
                    @click="restore(h)"
                  >Повторить</v-btn>
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
  points.value.push({ x: lastX + 1, y: 0 })
}

const removePoint = (index) => {
  points.value.splice(index, 1)
}

const calculate = async () => {
  if (points.value.length < 2) return

  loading.value = true
  error.value = null
  try {
    const resp = await axios.post(`${API_URL}/interpolate`, {
      method: method.value,
      points: points.value,
      target_x: targetX.value
    })
    result.value = resp.data.result
    curve.value = resp.data.curve
    cached.value = resp.data.cached
    fetchHistory()
  } catch (err) {
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
    history.value = resp.data
  } catch (err) {
    console.error('Failed to fetch history', err)
  }
}

const restore = (h) => {
  method.value = h.method
  points.value = JSON.parse(JSON.stringify(h.points))
  targetX.value = h.target_x
  calculate()
}

watch([method, points, targetX], () => {
  debouncedCalculate()
}, { deep: true })

onMounted(() => {
  calculate()
  fetchHistory()
})
</script>
