<template>
  <div class="not-found">
    <canvas ref="particlesCanvas" class="particles-canvas"></canvas>
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1 class="error-title">页面未找到</h1>
      <p class="error-description">
        抱歉，您访问的页面不存在或已被移除。
      </p>
      <div class="error-actions">
        <el-button type="primary" @click="goHome">
          返回首页
        </el-button>
        <el-button @click="goBack">
          返回上页
        </el-button>
      </div>
      
      <div class="suggestions">
        <h3>您可以尝试：</h3>
        <ul>
          <li>检查URL是否正确</li>
          <li>使用导航菜单浏览其他页面</li>
          <li>联系管理员获取帮助</li>
        </ul>
      </div>
    </div>
    
    <div class="not-found-illustration">
      <el-icon size="120" color="#e4e7ed"><QuestionFilled /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { onMounted, onBeforeUnmount, ref } from 'vue'

const router = useRouter()

// 返回首页
const goHome = () => {
  router.push('/')
}

// 返回上一页
const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}

// 进入 404 页面后 10 秒自动返回首页
let autoTimer: number | undefined
onMounted(() => {
  autoTimer = window.setTimeout(() => {
    router.push('/')
  }, 10000)

  // 初始化粒子效果
  initParticles()
  window.addEventListener('resize', resizeCanvas)
})

onBeforeUnmount(() => {
  if (autoTimer) {
    clearTimeout(autoTimer)
    autoTimer = undefined
  }
  cancelAnimationFrame(animationId)
  window.removeEventListener('resize', resizeCanvas)
})

// 粒子效果（与登录页一致风格的精简版）
const particlesCanvas = ref<HTMLCanvasElement | null>(null)
let ctx: CanvasRenderingContext2D | null = null
let animationId = 0

interface Particle { x: number; y: number; vx: number; vy: number }
let particles: Particle[] = []

const resizeCanvas = () => {
  if (!particlesCanvas.value) return
  const dpr = window.devicePixelRatio || 1
  const rect = particlesCanvas.value.getBoundingClientRect()
  particlesCanvas.value.width = Math.floor(rect.width * dpr)
  particlesCanvas.value.height = Math.floor(rect.height * dpr)
  ctx = particlesCanvas.value.getContext('2d')
  if (ctx) ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
}

const initParticles = () => {
  if (!particlesCanvas.value) return
  resizeCanvas()
  const rect = particlesCanvas.value.getBoundingClientRect()
  const area = rect.width * rect.height
  const density = 0.00012
  const count = Math.max(14, Math.min(200, Math.floor(area * density)))
  particles = Array.from({ length: count }).map(() => ({
    x: Math.random() * rect.width,
    y: Math.random() * rect.height,
    vx: (Math.random() - 0.5) * 0.6,
    vy: (Math.random() - 0.5) * 0.6
  }))
  animate()
}

const draw = () => {
  if (!ctx || !particlesCanvas.value) return
  const { width, height } = particlesCanvas.value.getBoundingClientRect()
  ctx.clearRect(0, 0, width, height)

  const maxDist = 120
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x
      const dy = particles[i].y - particles[j].y
      const dist = Math.hypot(dx, dy)
      if (dist < maxDist) {
        const alpha = 1 - dist / maxDist
        ctx.strokeStyle = `rgba(64,158,255,${alpha * 0.8})`
        ctx.lineWidth = 1
        ctx.beginPath()
        ctx.moveTo(particles[i].x, particles[i].y)
        ctx.lineTo(particles[j].x, particles[j].y)
        ctx.stroke()
      }
    }
  }

  for (const p of particles) {
    ctx.fillStyle = 'rgba(64,158,255,0.9)'
    ctx.beginPath()
    ctx.arc(p.x, p.y, 2, 0, Math.PI * 2)
    ctx.fill()
  }
}

const animate = () => {
  const { width, height } = particlesCanvas.value!.getBoundingClientRect()
  for (const p of particles) {
    p.x += p.vx
    p.y += p.vy
    if (p.x < 0 || p.x > width) p.vx *= -1
    if (p.y < 0 || p.y > height) p.vy *= -1
  }
  draw()
  animationId = requestAnimationFrame(animate)
}
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  min-width: 800px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(1200px 600px at 10% 10%, rgba(64, 158, 255, 0.12), transparent),
              radial-gradient(800px 400px at 90% 80%, rgba(103, 194, 58, 0.12), transparent),
              linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  padding: 20px;
}

/* 粒子画布，与登录页保持一致 */
.particles-canvas {
  position: fixed;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  opacity: 0.6;
  z-index: 0;
}

.not-found-content {
  text-align: center;
  max-width: 600px;
  min-width: 400px;
  margin-right: 60px;
  min-height: 400px;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  color: #409eff;
  line-height: 1;
  margin-bottom: 20px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.error-title {
  font-size: 32px;
  color: #303133;
  margin: 0 0 16px 0;
}

.error-description {
  font-size: 16px;
  color: #606266;
  margin: 0 0 32px 0;
  line-height: 1.6;
}

.error-actions {
  margin-bottom: 40px;
}

.error-actions .el-button {
  margin: 0 8px;
}

.suggestions {
  text-align: left;
  background: rgba(255, 255, 255, 0.8);
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.suggestions h3 {
  margin: 0 0 12px 0;
  color: #303133;
  font-size: 16px;
}

.suggestions ul {
  margin: 0;
  padding-left: 20px;
  color: #606266;
}

.suggestions li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.not-found-illustration {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .not-found {
    flex-direction: column;
    text-align: center;
  }
  
  .not-found-content {
    margin-right: 0;
    margin-bottom: 40px;
  }
  
  .error-code {
    font-size: 80px;
  }
  
  .error-title {
    font-size: 24px;
  }
  
  .suggestions {
    text-align: center;
  }
}
</style> 