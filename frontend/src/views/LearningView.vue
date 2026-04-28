<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="8" class="mx-auto">
        <v-card class="mb-6">
          <v-tabs v-model="tab" bg-color="primary">
            <v-tab value="theory">Теория</v-tab>
            <v-tab value="practice">Задания</v-tab>
          </v-tabs>

          <v-card-text>
            <v-window v-model="tab">
              <!-- Theory -->
              <v-window-item value="theory">
                <div v-for="t in theory" :key="t.title" class="mb-6">
                  <h2 class="text-h5 mb-2">{{ t.title }}</h2>
                  <p class="text-body-1">{{ t.content }}</p>
                  <v-divider class="mt-4"></v-divider>
                </div>
              </v-window-item>

              <!-- Practice -->
              <v-window-item value="practice">
                <v-alert type="info" class="mb-4">
                  Решите задачи, используя изученные методы. Допуск ошибки: 0.01.
                </v-alert>

                <div v-for="task in tasks" :key="task.id" class="mb-8">
                  <v-card variant="outlined" class="pa-4">
                    <v-card-title>{{ task.title }}</v-card-title>
                    <v-card-text>
                      <p class="mb-4">{{ task.description }}</p>
                      <v-row align="center">
                        <v-col cols="12" sm="6">
                          <v-text-field
                            v-model.number="userAnswers[task.id]"
                            label="Ваш ответ"
                            type="number"
                            :color="getAnswerColor(task.id)"
                            :append-inner-icon="getAnswerIcon(task.id)"
                            hide-details
                          ></v-text-field>
                        </v-col>
                        <v-col cols="12" sm="6">
                          <v-btn
                            color="primary"
                            @click="checkAnswer(task.id)"
                          >Проверить</v-btn>
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>
                </div>
              </v-window-item>
            </v-window>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import tasksData from '../tasks.json'

const tab = ref('theory')
const tasks = ref(tasksData)
const userAnswers = ref({})
const results = ref({})

const theory = [
  {
    title: 'Что такое интерполяция?',
    content: 'Интерполяция — способ нахождения промежуточных значений величины по имеющемуся дискретному набору известных значений. В обучении мы рассмотрим три метода: линейную, Лагранжа и Ньютона.'
  },
  {
    title: 'Линейная интерполяция',
    content: 'Самый простой вид интерполяции, при котором две соседние точки соединяются отрезком прямой. Формула: y = y0 + (x - x0) * (y1 - y0) / (x1 - x0).'
  },
  {
    title: 'Полином Лагранжа',
    content: 'Позволяет построить полином, проходящий через все заданные точки. Он не требует решения систем уравнений, но сложен в вычислении при добавлении новых точек.'
  },
  {
    title: 'Полином Ньютона',
    content: 'Использует разделенные разности. Его преимущество в том, что при добавлении новой точки достаточно вычислить только одно дополнительное слагаемое.'
  }
]

const checkAnswer = (id) => {
  const task = tasks.value.find(t => t.id === id)
  const userVal = userAnswers.value[id]
  if (userVal === undefined) return

  const isCorrect = Math.abs(userVal - task.correct_answer) < 0.01
  results.value[id] = isCorrect
}

const getAnswerColor = (id) => {
  if (results.value[id] === undefined) return ''
  return results.value[id] ? 'success' : 'error'
}

const getAnswerIcon = (id) => {
  if (results.value[id] === undefined) return ''
  return results.value[id] ? 'mdi-check-circle' : 'mdi-alert-circle'
}
</script>
