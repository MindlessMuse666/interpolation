<template>
  <div class="chart-container">
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
  CategoryScale
} from 'chart.js'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale
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
  targetPoint: {
    type: Object,
    default: null
  }
})

const chartData = computed(() => {
  const datasets = [
    {
      label: 'Исходные точки',
      data: props.points.map(p => ({ x: p.x, y: p.y })),
      backgroundColor: '#f87979',
      borderColor: '#f87979',
      showLine: false,
      pointRadius: 6,
      type: 'scatter'
    }
  ]

  if (props.curve.length > 0) {
    datasets.push({
      label: 'Интерполяционная кривая',
      data: props.curve.map(p => ({ x: p.x, y: p.y })),
      borderColor: '#3f51b5',
      backgroundColor: '#3f51b5',
      tension: 0.1,
      pointRadius: 0,
      fill: false
    })
  }

  if (props.targetPoint) {
    datasets.push({
      label: 'Результат',
      data: [{ x: props.targetPoint.x, y: props.targetPoint.y }],
      backgroundColor: '#4caf50',
      borderColor: '#4caf50',
      pointRadius: 8,
      pointStyle: 'star',
      type: 'scatter'
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
  height: 400px;
  width: 100%;
}
</style>
