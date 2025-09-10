<template>
  <div class="login-container">
    <canvas ref="particlesCanvas" class="particles-canvas"></canvas>
    <div class="login-box">
      <el-card class="login-card" shadow="hover">
        <div class="card-body">
          <div class="brand-pane">
            <div class="brand-header">
              <img
                class="brand-logo"
                src="../assets/async.png"
                alt="AsyncLab"
              />
              <h2>AsyncLab</h2>
              <p>欢迎每一个热爱开发的人</p>
            </div>
            <ul class="brand-highlights">
              <li>
                <el-icon><Promotion /></el-icon>
                <a href="https://www.asyncraft.club/" target="_blank"
                  >Asyncraft官网</a
                >
              </li>
              <li>
                <el-icon><Setting /></el-icon>
                <a href="https://github.com/async-lab" target="_blank"
                  >AsyncLab GitHub</a
                >
              </li>
              <li>
                <el-icon><Box /></el-icon>
                <a href="https://gitlab.asynclab.club/" target="_blank"
                  >AsyncLab Gitlab</a
                >
              </li>
            </ul>
          </div>

          <div class="form-pane">
            <div class="login-header">
              <h2>用户登录</h2>
              <p>使用您的账户继续</p>
            </div>

            <el-form
              ref="loginFormRef"
              :model="loginForm"
              :rules="loginRules"
              label-width="0"
              @submit.prevent="handleLogin"
              class="login-form"
            >
              <el-form-item prop="username">
                <el-input
                  v-model.trim="loginForm.username"
                  placeholder="用户名"
                  :prefix-icon="User"
                  clearable
                  size="large"
                  @keyup.enter="handleLogin"
                />
              </el-form-item>

              <el-form-item prop="password">
                <el-input
                  v-model.trim="loginForm.password"
                  type="password"
                  placeholder="密码"
                  :prefix-icon="Lock"
                  show-password
                  clearable
                  size="large"
                  @keyup.enter="handleLogin"
                />
              </el-form-item>

              <div class="form-extras">
                <el-checkbox v-model="remember">记住用户名</el-checkbox>
                <el-checkbox v-model="rememberPassword">记住密码</el-checkbox>
              </div>

              <el-button
                type="primary"
                :loading="loading"
                :disabled="!canSubmit || loading"
                @click="handleLogin"
                size="large"
                class="submit-btn"
              >
                {{ loading ? "登录中..." : "登录" }}
              </el-button>

              <div class="footer-links">
                <el-button type="text" class="link" @click="goToHome"
                  >返回首页</el-button
                >
              </div>
            </el-form>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import type { FormInstance, FormRules } from "element-plus";
import {
  getUsername,
  setUsername,
  setToken,
  setUserProfile,
} from "@/utils/auth";
import { getMeInfo } from "@/api/user";
import {
  saveEncryptedPassword,
  loadEncryptedPassword,
  clearEncryptedPassword,
} from "@/utils/auth";
import { useFailedTip, useSuccessTip } from "@/utils/msgTip";
import { createToken } from "@/api/auth";
import type { LoginRequest } from "@/api/types";
import { Box, Promotion, Setting } from "@element-plus/icons-vue";
import { User, Lock } from "@element-plus/icons-vue";

const router = useRouter();
const route = useRoute();
const loginFormRef = ref<FormInstance>();
const loading = ref(false);

// 登录表单数据
const loginForm = reactive<LoginRequest>({
  username: "",
  password: "",
});

// 记住用户名选项
const remember = ref(false);
const rememberPassword = ref(false);

// 可提交状态
const canSubmit = computed(() => {
  return (
    loginForm.username.trim().length > 0 && loginForm.password.trim().length > 0
  );
});

// 表单验证规则
const loginRules: FormRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return;
  if (!canSubmit.value || loading.value) return;

  try {
    await loginFormRef.value.validate();
    loading.value = true;

    // 调用登录API
    const tokenResp = (await createToken({
      username: loginForm.username.trim(),
      password: loginForm.password.trim(),
    })) as any;

    // 解析token字符串
    const token = tokenResp.data;

    // 登录成功，保存token
    if (token) {
      setToken(token);

      // 记住用户名
      if (remember.value) {
        setUsername(loginForm.username.trim());
      } else {
        setUsername("");
      }

      // 记住密码（加密存储）
      if (rememberPassword.value) {
        await saveEncryptedPassword(
          loginForm.username.trim(),
          loginForm.password.trim()
        );
      } else {
        clearEncryptedPassword();
      }

      // 获取并保存当前用户信息（完整对象）
      try {
        const me = (await getMeInfo()).data as any;
        if (me) {
          setUserProfile(me);
        }
      } catch (e) {
        // 忽略获取用户信息失败
      }

      useSuccessTip("登录成功");

      // 跳转到目标页面或首页
      const redirect = route.query.redirect as string;
      router.push(redirect || "/dashboard");
    } else {
      useFailedTip("登录失败，请检查用户名和密码");
    }
  } catch (error: any) {
    console.error("登录失败:", error);
    useFailedTip(
      error?.msg || error?.message || "登录失败，请检查用户名和密码"
    );
  } finally {
    loading.value = false;
  }
};

// 返回首页
const goToHome = () => {
  router.push("/");
};

// 组件挂载时加载记住的用户名
onMounted(async () => {
  const savedUsername = getUsername();
  if (savedUsername) {
    loginForm.username = savedUsername;
    remember.value = true;
  }
  // 尝试加载加密密码
  if (savedUsername) {
    const savedPwd = await loadEncryptedPassword(savedUsername);
    if (savedPwd) {
      loginForm.password = savedPwd;
      rememberPassword.value = true;
    }
  }
  initParticles();
  window.addEventListener("resize", resizeCanvas);
});

onBeforeUnmount(() => {
  cancelAnimationFrame(animationId);
  window.removeEventListener("resize", resizeCanvas);
});

// 粒子效果
const particlesCanvas = ref<HTMLCanvasElement | null>(null);
let ctx: CanvasRenderingContext2D | null = null;
let animationId = 0;

interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
}
let particles: Particle[] = [];

const resizeCanvas = () => {
  if (!particlesCanvas.value) return;
  const dpr = window.devicePixelRatio || 1;
  const rect = particlesCanvas.value.getBoundingClientRect();
  particlesCanvas.value.width = Math.floor(rect.width * dpr);
  particlesCanvas.value.height = Math.floor(rect.height * dpr);
  ctx = particlesCanvas.value.getContext("2d");
  if (ctx) ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
};

const initParticles = () => {
  if (!particlesCanvas.value) return;
  resizeCanvas();
  const rect = particlesCanvas.value.getBoundingClientRect();
  const area = rect.width * rect.height;
  const density = 0.00012; // 粒子密度（进一步降低）
  const count = Math.max(14, Math.min(200, Math.floor(area * density)));
  particles = Array.from({ length: count }).map(() => ({
    x: Math.random() * rect.width,
    y: Math.random() * rect.height,
    vx: (Math.random() - 0.5) * 0.6,
    vy: (Math.random() - 0.5) * 0.6,
  }));
  animate();
};

const draw = () => {
  if (!ctx || !particlesCanvas.value) return;
  const { width, height } = particlesCanvas.value.getBoundingClientRect();
  ctx.clearRect(0, 0, width, height);

  // 连线
  const maxDist = 120;
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x;
      const dy = particles[i].y - particles[j].y;
      const dist = Math.hypot(dx, dy);
      if (dist < maxDist) {
        const alpha = 1 - dist / maxDist;
        ctx.strokeStyle = `rgba(64,158,255,${alpha * 0.8})`;
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(particles[i].x, particles[i].y);
        ctx.lineTo(particles[j].x, particles[j].y);
        ctx.stroke();
      }
    }
  }

  // 粒子
  for (const p of particles) {
    ctx.fillStyle = "rgba(64,158,255,0.9)";
    ctx.beginPath();
    ctx.arc(p.x, p.y, 2, 0, Math.PI * 2);
    ctx.fill();
  }
};

const animate = () => {
  const { width, height } = particlesCanvas.value!.getBoundingClientRect();
  for (const p of particles) {
    p.x += p.vx;
    p.y += p.vy;
    if (p.x < 0 || p.x > width) p.vx *= -1;
    if (p.y < 0 || p.y > height) p.vy *= -1;
  }
  draw();
  animationId = requestAnimationFrame(animate);
};
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(
      1200px 600px at 10% 10%,
      rgba(64, 158, 255, 0.12),
      transparent
    ),
    radial-gradient(
      800px 400px at 90% 80%,
      rgba(103, 194, 58, 0.12),
      transparent
    ),
    linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  position: relative;
}

/* 左下角粒子画布 */
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

.login-box {
  width: 100%;
  max-width: 900px;
  position: relative;
  z-index: 1;
}

.login-card {
  border-radius: 16px;
  overflow: hidden;
}

.card-body {
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 0;
}

.brand-pane {
  background: linear-gradient(180deg, #ecf5ff 0%, #ffffff 100%);
  padding: 40px 32px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
}

.brand-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.brand-logo {
  width: 56px;
  height: 56px;
}

.brand-header h2 {
  margin: 0;
  color: #303133;
}

.brand-header p {
  margin: 0;
  color: #606266;
}

.brand-highlights {
  list-style: none;
  padding: 0;
  margin: 16px 0 0 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.brand-highlights li {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #409eff;
}

.brand-highlights a {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
}

.brand-highlights a::before {
  content: "";
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(64, 158, 255, 0.1),
    transparent
  );
  transition: left 0.5s ease;
}

.brand-highlights a:hover {
  color: #ffffff;
  background: linear-gradient(135deg, #409eff 0%, #66b3ff 100%);
  border-color: #409eff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.brand-highlights a:hover::before {
  left: 100%;
}

.brand-highlights a:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.form-pane {
  padding: 40px 32px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login-header {
  text-align: left;
  margin-bottom: 24px;
}

.login-header h2 {
  color: #303133;
  margin: 0 0 8px 0;
  font-size: 24px;
}

.login-header p {
  color: #909399;
  margin: 0;
  font-size: 14px;
}

.login-form {
  width: 100%;
}

.form-extras {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.link {
  padding: 0;
}

.submit-btn {
  width: 100%;
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
  margin-top: 12px;
}

@media (max-width: 1024px) {
  .login-box {
    max-width: 720px;
  }
  .card-body {
    grid-template-columns: 1fr;
  }
  .brand-pane {
    display: none;
  }
  .particles-canvas {
    width: 100vw;
    height: 100vh;
  }
}

/* 更小屏幕的自适应（手机） */
@media (max-width: 768px) {
  .login-container {
    padding: 16px;
  }
  .login-box {
    max-width: 520px;
  }
  .form-pane {
    padding: 24px 20px;
  }
}
</style>
<style>
/* 暗色模式覆盖：Login 背景与品牌面板 */
html.dark .login-container {
  background: radial-gradient(
      1200px 600px at 10% 10%,
      rgba(96, 165, 250, 0.10),
      transparent
    ),
    radial-gradient(
      800px 400px at 90% 80%,
      rgba(134, 239, 172, 0.10),
      transparent
    ),
    /* 星云彩色雾气层 */
    radial-gradient(900px 600px at 20% 30%, rgba(56, 189, 248, 0.14), transparent 60%),
    radial-gradient(1000px 700px at 78% 68%, rgba(168, 85, 247, 0.12), transparent 60%),
    radial-gradient(800px 500px at 40% 80%, rgba(59, 130, 246, 0.10), transparent 60%),
    linear-gradient(135deg, #0b0f14 0%, #0f1720 100%);
  animation: space-drift 40s linear infinite;
}

html.dark .login-card {
  background: var(--card-bg);
}

html.dark .brand-pane {
  background: linear-gradient(180deg, #0f1720 0%, #0b0f14 100%);
}

html.dark .brand-header h2 { color: var(--el-text-color-primary); }
html.dark .brand-header p { color: var(--el-text-color-regular); }
html.dark .brand-highlights li { color: var(--el-color-primary); }
html.dark .brand-highlights a { color: var(--el-color-primary); border-color: transparent; }
html.dark .brand-highlights a:hover { color: #ffffff; }

html.dark .login-header h2 { color: var(--el-text-color-primary); }
html.dark .login-header p { color: var(--el-text-color-regular); }

/* 星空背景与关闭粒子（暗色模式） */
html.dark .particles-canvas { display: none; }
html.dark .login-container::before {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background:
    radial-gradient(1px 1px at 20% 30%, rgba(255,255,255,0.98) 50%, transparent 51%),
    radial-gradient(1px 1px at 40% 70%, rgba(255,255,255,0.95) 50%, transparent 51%),
    radial-gradient(1px 1px at 65% 25%, rgba(255,255,255,0.9) 50%, transparent 51%),
    radial-gradient(1px 1px at 80% 60%, rgba(255,255,255,0.98) 50%, transparent 51%),
    radial-gradient(1px 1px at 15% 85%, rgba(255,255,255,0.88) 50%, transparent 51%),
    radial-gradient(1px 1px at 55% 55%, rgba(255,255,255,0.95) 50%, transparent 51%),
    radial-gradient(1px 1px at 30% 50%, rgba(255,255,255,0.9) 50%, transparent 51%),
    radial-gradient(1px 1px at 72% 40%, rgba(255,255,255,0.92) 50%, transparent 51%),
    radial-gradient(1px 1px at 88% 20%, rgba(255,255,255,0.95) 50%, transparent 51%),
    radial-gradient(1px 1px at 8% 64%, rgba(255,255,255,0.85) 50%, transparent 51%),
    radial-gradient(1px 1px at 12% 52%, rgba(255,255,255,0.9) 50%, transparent 51%),
    radial-gradient(1px 1px at 48% 82%, rgba(255,255,255,0.92) 50%, transparent 51%),
    radial-gradient(1px 1px at 68% 12%, rgba(255,255,255,0.9) 50%, transparent 51%),
    radial-gradient(1px 1px at 92% 48%, rgba(255,255,255,0.95) 50%, transparent 51%);
  background-size: auto;
  animation: twinkle 2.2s infinite ease-in-out, star-drift-1 60s linear infinite;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.18; }
  50% { opacity: 1; }
}
@keyframes star-drift-1 {
  0% { background-position: 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px; }
  100% { background-position: 80px 60px, -60px 40px, 100px -40px, -80px -60px, 40px -80px, -100px 100px, 60px -60px, -40px 80px, 120px 20px, -120px -20px, 90px -90px, -90px 90px, 70px 40px, -70px -40px; }
}

/* 第二层较大的星点，错相闪烁，提升可见度 */
html.dark .login-container::after {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background:
    radial-gradient(1.5px 1.5px at 12% 22%, rgba(255,255,255,0.98) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 52% 18%, rgba(255,255,255,0.92) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 78% 46%, rgba(255,255,255,0.98) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 34% 82%, rgba(255,255,255,0.9) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 90% 72%, rgba(255,255,255,1) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 26% 36%, rgba(255,255,255,0.95) 50%, transparent 51%),
    radial-gradient(1.5px 1.5px at 68% 68%, rgba(255,255,255,0.98) 50%, transparent 51%);
  background-size: auto;
  animation: twinkle2 2.8s infinite ease-in-out alternate, star-drift-2 90s linear infinite;
}

@keyframes twinkle2 {
  0% { opacity: 0.15; }
  50% { opacity: 1; }
  100% { opacity: 0.3; }
}
@keyframes star-drift-2 {
  0% { background-position: 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px, 0px 0px; }
  100% { background-position: -120px 80px, 100px -60px, -80px -100px, 60px 120px, -140px 40px, 90px -90px, -60px 60px; }
}

@keyframes space-drift {
  0% { background-position: 0 0, 0 0, 0 0, 0 0, 0 0, 0 0, 0 0; }
  100% { background-position: 40px 60px, -60px 40px, 30px -30px, -40px -60px, 20px -20px, -30px 30px, 0 0; }
}
</style>
