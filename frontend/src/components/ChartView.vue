<template>
  <div class="chart-container" :style="containerStyle">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  LineController,
  ScatterController
} from 'chart.js'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  LineController,
  ScatterController
)

const props = defineProps({
  points: {
    type: Array,
    required: true
  },
  curve: {
    type: Array,
    default: () => []
  },
  height: {
    type: [Number, String],
    default: 400
  },
  targetPoint: {
    type: Object,
    default: null
  }
})

const containerStyle = computed(() => {
  if (typeof props.height === 'number') return { height: `${props.height}px` }
  if (typeof props.height === 'string' && props.height.trim().length > 0) return { height: props.height }
  return { height: '400px' }
})

const chartData = computed(() => {
  const datasets = [
    {
      label: 'Исходные точки',
      data: props.points.map(p => ({ x: Number(p.x), y: Number(p.y) })),
      backgroundColor: '#E88CA5',
      borderColor: '#E88CA5',
      showLine: false,
      pointRadius: 6,
      type: 'scatter',
      order: 1
    }
  ]

  if (props.curve && props.curve.length > 0) {
    datasets.push({
      label: 'Интерполяционная кривая',
      data: props.curve.map(p => ({ x: Number(p.x), y: Number(p.y) })),
      borderColor: '#5E9DC8',
      backgroundColor: '#5E9DC8',
      tension: 0.3,
      pointRadius: 0,
      fill: false,
      borderWidth: 3,
      order: 2
    })
  }

  if (props.targetPoint) {
    datasets.push({
      label: 'Результат',
      data: [{ x: Number(props.targetPoint.x), y: Number(props.targetPoint.y) }],
      backgroundColor: '#3BAB7B',
      borderColor: '#3BAB7B',
      pointRadius: 10,
      pointStyle: 'rectRot',
      type: 'scatter',
      order: 0
    })
  }

  return { datasets }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      type: 'linear',
      position: 'bottom'
    },
    y: {
      type: 'linear'
    }
  }
}
</script>

<style scoped>
.chart-container {
  width: 100%;
}
</style>
