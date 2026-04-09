<template>
  <div class="app-wrapper">
    <!-- Navbar -->
    <nav class="tf-navbar">
      <NuxtLink to="/" class="tf-brand">
        <div class="brand-icon">
            <i class="fa-solid fa-book"></i>        </div>
        <span class="brand-name">Task<span>Manager</span></span>
      </NuxtLink>

      <div class="tf-nav-stats" v-if="stats">
        <div class="stat-item">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">Total</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">Pendientes</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.completed }}</div>
          <div class="stat-label">Completadas</div>
        </div>
      </div>
    </nav>

    <!-- Slot -->
    <main class="main-content">
      <div class="container">
        <slot />
      </div>
    </main>

    <!-- Toast container -->
    <ToastContainer />
  </div>
</template>

<script setup>
import { useTasksStore } from "~/store/tasks";

const tasksStore = useTasksStore();
const stats = computed(() => tasksStore.stats);
</script>