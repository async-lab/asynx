<template>
  <section class="home-hero" :class="compact ? 'compact' : ''">
    <canvas ref="particlesCanvas" class="hero-particles"></canvas>
    <div class="hero-main">
      <img class="hero-logo" src="../assets/async.png" alt="AsyncLab" />
      <div class="hero-texts">
        <h1 class="hero-title">{{ title }}</h1>
        <p v-if="subtitle" class="hero-subtitle">{{ subtitle }}</p>
      </div>
    </div>

    <nav class="hero-links" v-if="resolvedLinks.length">
      <a
        v-for="link in resolvedLinks"
        :key="link.href"
        class="hero-link"
        :href="link.href"
        target="_blank"
        rel="noopener"
      >
        {{ link.label }}
      </a>
    </nav>

    <div class="hero-extra">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onBeforeUnmount } from "vue";

interface LinkItem {
  label: string;
  href: string;
}

const props = withDefaults(
  defineProps<{
    title?: string;
    subtitle?: string;
    links?: LinkItem[];
    compact?: boolean;
  }>(),
  {
    title: "AsyncLab",
    subtitle: "欢迎每一个热爱开发的人",
    links: () => [
      { label: "Asyncraft 官网", href: "https://www.asyncraft.club/" },
      { label: "AsyncLab GitHub", href: "https://github.com/async-lab" },
      { label: "AsyncLab Gitlab", href: "https://gitlab.asynclab.club/" },
    ],
    compact: false,
  }
);

const resolvedLinks = computed(() =>
  Array.isArray(props.links) ? props.links : []
);

// 粒子背景
const particlesCanvas = ref<HTMLCanvasElement | null>(null);
let animationFrameId = 0;
let cleanupResize: (() => void) | null = null;

interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  radius: number;
}

function createParticles(
  count: number,
  width: number,
  height: number
): Particle[] {
  const particles: Particle[] = [];
  for (let i = 0; i < count; i++) {
    particles.push({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * 0.6,
      vy: (Math.random() - 0.5) * 0.6,
      radius: 1 + Math.random() * 1.8,
    });
  }
  return particles;
}

function startParticles() {
  const canvasMaybe = particlesCanvas.value;
  if (!canvasMaybe) return;
  const canvas: HTMLCanvasElement = canvasMaybe;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;
  const context = ctx as CanvasRenderingContext2D;

  const parent = canvas.parentElement as HTMLElement;
  const dpi = Math.min(window.devicePixelRatio || 1, 2);

  function resize() {
    const rect = parent.getBoundingClientRect();
    canvas.width = Math.max(1, Math.floor(rect.width * dpi));
    canvas.height = Math.max(1, Math.floor(rect.height * dpi));
    canvas.style.width = rect.width + "px";
    canvas.style.height = rect.height + "px";
  }

  resize();

  const area = () => (canvas.width * canvas.height) / (dpi * dpi);
  const baseCount = Math.max(24, Math.min(90, Math.floor(area() / 12000)));
  let particles = createParticles(baseCount, canvas.width, canvas.height);

  function tick() {
    context.clearRect(0, 0, canvas.width, canvas.height);

    for (let i = 0; i < particles.length; i++) {
      const p = particles[i];
      p.x += p.vx;
      p.y += p.vy;
      if (p.x < 0 || p.x > canvas.width) p.vx *= -1;
      if (p.y < 0 || p.y > canvas.height) p.vy *= -1;
    }

    context.strokeStyle = "rgba(46,110,233,0.15)";
    context.lineWidth = 1 * dpi;
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist2 = dx * dx + dy * dy;
        const maxDist = 120 * dpi;
        if (dist2 < maxDist * maxDist) {
          const opacity = 1 - Math.sqrt(dist2) / maxDist;
          context.globalAlpha = Math.max(0.05, Math.min(0.5, opacity * 0.6));
          context.beginPath();
          context.moveTo(particles[i].x, particles[i].y);
          context.lineTo(particles[j].x, particles[j].y);
          context.stroke();
        }
      }
    }
    context.globalAlpha = 1;

    context.fillStyle = "rgba(46,110,233,0.55)";
    for (const p of particles) {
      context.beginPath();
      context.arc(p.x, p.y, p.radius * dpi, 0, Math.PI * 2);
      context.fill();
    }

    animationFrameId = requestAnimationFrame(tick);
  }

  animationFrameId = requestAnimationFrame(tick);

  const onResize = () => {
    const oldW = canvas.width;
    const oldH = canvas.height;
    resize();
    const scaleX = canvas.width / (oldW || 1);
    const scaleY = canvas.height / (oldH || 1);
    particles = particles.map((p) => ({
      ...p,
      x: p.x * scaleX,
      y: p.y * scaleY,
    }));
  };
  window.addEventListener("resize", onResize);
  cleanupResize = () => window.removeEventListener("resize", onResize);
}

onMounted(() => {
  startParticles();
});

onBeforeUnmount(() => {
  if (animationFrameId) cancelAnimationFrame(animationFrameId);
  if (cleanupResize) cleanupResize();
});
</script>

<style scoped>
.home-hero {
  width: 98%;
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 28px;
  border-radius: 18px;
  height: 500px;
  position: relative;
  overflow: hidden;
  background: linear-gradient(
      180deg,
      rgba(255, 255, 255, 0.9),
      rgba(255, 255, 255, 0.85)
    ),
    radial-gradient(
      900px 600px at 8% 12%,
      rgba(64, 158, 255, 0.12),
      transparent 60%
    ),
    radial-gradient(
      700px 420px at 92% 80%,
      rgba(103, 194, 58, 0.12),
      transparent 60%
    ),
    linear-gradient(135deg, #f6f9ff 0%, #ffffff 60%);
  box-shadow: 0 10px 30px rgba(31, 45, 61, 0.1);
}

/* 细网格背景（装饰） */
.home-hero::before {
  content: "";
  position: absolute;
  inset: 0;
  background-image: linear-gradient(
      to right,
      rgba(64, 158, 255, 0.08) 1px,
      transparent 1px
    ),
    linear-gradient(to bottom, rgba(64, 158, 255, 0.08) 1px, transparent 1px);
  background-size: 22px 22px;
  mask-image: radial-gradient(
    600px 400px at 20% 20%,
    rgba(0, 0, 0, 0.6),
    transparent 70%
  );
  pointer-events: none;
  z-index: 0;
}

.home-hero.compact {
  padding: 18px 20px;
  gap: 12px;
  border-radius: 14px;
}

.hero-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 360px;
  position: relative;
  z-index: 1;
}

.hero-logo {
  width: 120px;
  height: 120px;
  border-radius: 18px;
}
.hero-texts {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  position: relative;
  z-index: 1;
}
.hero-title {
  margin: 0;
  color: #1f2d3d;
  font-size: 36px;
  line-height: 1.15;
  letter-spacing: 0.2px;
}
.hero-subtitle {
  margin: 8px 0 0 0;
  color: #5e6d82;
  font-size: 16px;
}

.hero-links {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
  z-index: 1;
}
.hero-link {
  color: #2e6ee9;
  text-decoration: none;
  padding: 8px 14px;
  border-radius: 999px;
  border: 1px solid rgba(46, 110, 233, 0.16);
  background: rgba(46, 110, 233, 0.06);
  transition: all 0.2s ease;
  font-weight: 600;
}
.hero-link:hover {
  color: #fff;
  background: linear-gradient(135deg, #3a81ff, #66b3ff);
  border-color: transparent;
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(58, 129, 255, 0.25);
}

.hero-extra {
  margin-top: 4px;
}

.hero-particles {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

@media (max-width: 1200px) {
  .hero-title {
    font-size: 30px;
  }
  .hero-logo {
    width: 96px;
    height: 96px;
  }
}

@media (max-width: 768px) {
  .home-hero {
    padding: 16px;
    border-radius: 12px;
    gap: 10px;
  }
  .hero-title {
    font-size: 24px;
  }
  .hero-subtitle {
    font-size: 14px;
  }
  .hero-logo {
    width: 72px;
    height: 72px;
  }
}
</style>
