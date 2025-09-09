<template>
  <div class="home">
    <el-card class="welcome-card">
      <template #header>
        <div class="card-header">
          <h1>欢迎来到 AsyncLab</h1>
        </div>
      </template>

      <div class="welcome-content">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="24" :md="12">
            <div class="welcome-text">
              <h2>异步开发实验室</h2>
              <div>
                <p>欢迎每一个热爱开发的人！</p>
                <p>
                  本网站为异步开发实验室的官方网站，旨在学习交流开发经验，共同进步。
                </p>
                <p>欢迎加入我们，一起探索开发的世界！</p>
                <p>
                  在这里你可以探索各种开发技术，学习各种开发知识，提高自己的开发能力。
                </p>
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="12">
            <div class="welcome-actions">
              <img
                class="brand-logo"
                src="../assets/async.png"
                alt="AsyncLab"
              />
              <el-button
                type="primary"
                size="large"
                @click="goToLogin"
                v-if="!isLoggedIn"
              >
                立即登录
              </el-button>
              <el-button
                type="success"
                size="large"
                @click="goToDashboard"
                v-if="isLoggedIn"
              >
                进入仪表板
              </el-button>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card class="features-card">
      <template #header>
        <div class="card-header">
          <h3>AsyncLab 加入我们</h3>
        </div>
      </template>

      <div class="join-content">
        <div class="poster-container">
          <div class="intro-text">
            <p>
              异步开发实验室全称
              <strong>Asynchronous Lab</strong>，隶属于网络空间安全学院（前身为
              IOTU），主要开设前端、后端（Java、Golang）、移动端（Android）、嵌入式等技术方向。
            </p>

            <p>
              异步，可能有的同学还不太了解它的含义。简言之，可以理解为两个人向他们各自的目标迈进，但是互相不用等着对方，只管自己有条不紊前进就好。我们不需要每个人的进度同步，只是希望同学们通过加入实验室，都能离自己的梦想更进一步，找到自己所爱，这也是实验室的意义所在。
            </p>

            <p>
              我们实验室成立的初衷就是为了给各位喜欢开发、热爱创造的同学们，提供良好的学习氛围、安静的学习场所，以及与有着同样爱好的同学交流、共同成长的机会。
            </p>

            <p>
              另外，实验室于 2013 年成立，到现在已经有 9
              年的历史了，所以会有很多已经毕业工作的学长给还在学校的我们提供很多帮助，比如互联网大厂内推、业界一线的情况等等。
            </p>
          </div>
          <img
            class="join-poster"
            src="../assets/海报.png"
            alt="AsyncLab 加入我们"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { getToken } from "@/utils/auth";
// 恢复到原页面结构，不再使用 HomeHero 组件

const router = useRouter();

// 计算是否已登录
const isLoggedIn = computed(() => {
  return !!getToken();
});

// 跳转到登录页
const goToLogin = () => {
  router.push("/login");
};

// 跳转到仪表板
const goToDashboard = () => {
  router.push("/dashboard");
};
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px;
  min-height: 600px;
}

.welcome-card {
  margin-bottom: 20px;
  min-height: 200px;
}

.features-card {
  margin-bottom: 20px;
  min-height: 200px;
}

.content-card {
  margin-bottom: 20px;
  min-height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 50px;
}

.card-header h1 {
  margin: 0;
  color: #409eff;
  font-size: 24px;
}

.card-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
}

.welcome-content {
  padding: 16px 0;
  min-height: 150px;
}

.welcome-text h2 {
  color: #303133;
  margin-bottom: 16px;
}

.welcome-text p {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 12px;
}

.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 120px;
}

.dynamic-content {
  padding: 20px 0;
  min-height: 300px;
}

.content-title h2 {
  color: #303133;
  margin-bottom: 16px;
  text-align: center;
}

.content-description p {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 24px;
  text-align: center;
  font-size: 16px;
}

.content-images {
  margin-bottom: 24px;
}

.images-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  align-items: center;
}

.images-container.vertical {
  flex-direction: column;
  max-width: 600px;
  margin: 0 auto;
}

.images-container.horizontal {
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
}

.image-wrapper {
  display: flex;
  justify-content: center;
}

.content-image {
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  object-fit: cover;
}

.content-features {
  margin-bottom: 24px;
}

.feature-card {
  text-align: center;
  padding: 20px;
  height: 100%;
}

.feature-icon {
  margin-bottom: 16px;
}

.feature-card h4 {
  color: #303133;
  margin-bottom: 12px;
  font-size: 16px;
}

.feature-card p {
  color: #606266;
  line-height: 1.5;
  margin: 0;
  font-size: 14px;
}
.brand-logo {
  width: 180px;
  height: 180px;
}

.join-content {
  padding: 20px 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.poster-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  max-width: 800px;
  margin: 0 auto;
}

@media (max-width: 1200px) {
  .home {
    padding: 16px;
  }
}

@media (max-width: 992px) {
  .brand-logo {
    width: 160px;
    height: 160px;
  }
}

@media (max-width: 768px) {
  .home {
    padding: 12px;
  }
  .welcome-content {
    padding: 12px 0;
  }
  .brand-logo {
    width: 140px;
    height: 140px;
  }
}

@media (max-width: 576px) {
  .brand-logo {
    width: 120px;
    height: 120px;
  }
}

.join-poster {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.join-poster:hover {
  transform: scale(1.02);
}

.intro-text {
  text-align: left;
  line-height: 1.8;
  color: #606266;
}

.intro-text p {
  margin-bottom: 16px;
  font-size: 15px;
  text-indent: 2em;
}

.intro-text p:last-child {
  margin-bottom: 0;
}

.intro-text strong {
  color: #409eff;
  font-weight: 600;
}
</style>
