<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'

const router = useRouter()

type Step = 'login' | 'totp' | 'totp-setup'
const step = ref<Step>('login')
const loading = ref(false)

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const rememberMe = ref(false)
const code = ref('')
const otp = ref<string[]>(['', '', '', '', '', ''])
const otpInputs = ref<HTMLInputElement[]>([])
const totpCode = computed(() => otp.value.join(''))

watch(step, (s) => {
  if (s === 'totp') nextTick(() => otpInputs.value[0]?.focus())
})
const tempToken = ref('')
const tempRememberMe = ref(false)
const totpSecret = ref('')
const totpQrUrl = ref('')

function storeToken(token: string) {
  localStorage.setItem('access_token', token)
}

async function fetchMe() {
  try {
    const { data } = await api.get('/auth/me')
    if (data?.id) localStorage.setItem('user_id', String(data.id))
  } catch {
    /* non-fatal */
  }
  try {
    const { data } = await api.get('/users')
    localStorage.setItem('is_super', data?.is_super ? '1' : '0')
  } catch {
    localStorage.setItem('is_super', '0')
  }
}

async function storeTokenAndRedirect(token: string) {
  localStorage.setItem('access_token', token)
  await fetchMe()
  router.push('/')
}

async function submitLogin() {
  if (!email.value || !password.value) {
    toast.warning('Email and password are required')
    return
  }
  loading.value = true
  try {
    const { data } = await api.post('/auth/login', {
      email: email.value,
      password: password.value,
      remember_me: rememberMe.value,
    })
    if (data.requires_totp) {
      tempToken.value = data.temp_token
      tempRememberMe.value = rememberMe.value
      step.value = 'totp'
    } else if (data.totp_setup_required) {
      storeToken(data.access_token)
      await setupTotp()
    } else {
      storeTokenAndRedirect(data.access_token)
    }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Login failed')
  } finally {
    loading.value = false
  }
}

async function setupTotp() {
  loading.value = true
  try {
    const { data } = await api.post('/auth/totp/setup')
    totpSecret.value = data.secret
    totpQrUrl.value = data.qr_url
    step.value = 'totp-setup'
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'TOTP setup failed')
  } finally {
    loading.value = false
  }
}

function onOtpInput(e: Event, i: number) {
  const v = (e.target as HTMLInputElement).value.replace(/\D/g, '')
  if (!v) {
    otp.value[i] = ''
    return
  }
  otp.value[i] = v[v.length - 1]
  if (i < 5) otpInputs.value[i + 1]?.focus()
  if (totpCode.value.length === 6) submitTotp()
}

function onOtpKeydown(e: KeyboardEvent, i: number) {
  if (e.key === 'Backspace' && !otp.value[i] && i > 0) {
    otp.value[i - 1] = ''
    otpInputs.value[i - 1]?.focus()
    e.preventDefault()
  } else if (e.key === 'ArrowLeft' && i > 0) {
    otpInputs.value[i - 1]?.focus()
  } else if (e.key === 'ArrowRight' && i < 5) {
    otpInputs.value[i + 1]?.focus()
  }
}

function onOtpPaste(e: ClipboardEvent, i: number) {
  e.preventDefault()
  const pasted = (e.clipboardData?.getData('text') || '').replace(/\D/g, '').slice(0, 6 - i)
  if (!pasted) return
  pasted.split('').forEach((c, idx) => {
    otp.value[i + idx] = c
  })
  const focusIdx = Math.min(i + pasted.length, 5)
  otpInputs.value[focusIdx]?.focus()
  if (totpCode.value.length === 6) submitTotp()
}

async function submitTotp() {
  if (totpCode.value.length !== 6) {
    toast.warning('Enter the 6-digit code')
    return
  }
  loading.value = true
  try {
    const { data } = await api.post('/auth/totp', {
      temp_token: tempToken.value,
      code: totpCode.value,
    })
    storeTokenAndRedirect(data.access_token)
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Invalid code')
  } finally {
    loading.value = false
  }
}

async function submitTotpSetup() {
  if (!code.value) {
    toast.warning('Code is required')
    return
  }
  loading.value = true
  try {
    const { data } = await api.post('/auth/totp/verify-setup', { code: code.value })
    if (data.access_token) storeToken(data.access_token)
    toast.success('TOTP verified')
    router.push('/')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Invalid code')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="grid min-h-screen lg:grid-cols-[1.1fr_1fr]">
    <!-- Brand panel -->
    <aside
      class="relative hidden flex-col justify-between overflow-hidden bg-primary p-12 text-primary-foreground lg:flex"
    >
      <div class="relative z-10 flex items-center gap-2">
        <span class="font-heading text-2xl font-bold tracking-tight">okuru.id</span>
      </div>

      <div class="relative z-10 max-w-md">
        <p class="font-heading text-4xl leading-tight font-bold tracking-tight">
          Mengalirkan
          <span class="text-white/70">nilai</span>.
        </p>
        <p class="mt-4 text-sm leading-relaxed text-primary-foreground/70">
          Uang, data, aset, dan informasi — bergerak secara aman, efisien, dan bermakna.
          Terinspirasi dari <span class="font-mono">okuru</span> (送る / 贈る).
        </p>
      </div>

      <div class="relative z-10 text-xs text-primary-foreground/50">
        &copy; {{ new Date().getFullYear() }} okuru.id
      </div>
    </aside>

    <!-- Form panel -->
    <div class="flex min-h-screen items-center justify-center bg-background p-6 sm:p-12">
      <div class="w-full max-w-sm">
        <!-- mobile wordmark -->
        <div class="mb-8 flex items-center justify-center lg:hidden">
          <span class="font-heading text-2xl font-bold tracking-tight">okuru.id</span>
        </div>

        <div class="mb-8">
          <h1 class="font-heading text-2xl font-bold tracking-tight">
            <template v-if="step === 'login'">Sign in</template>
            <template v-else-if="step === 'totp'">Two-factor</template>
            <template v-else>Set up 2FA</template>
          </h1>
          <p class="mt-1.5 text-sm text-muted-foreground">
            <template v-if="step === 'login'">Welcome back. Sign in to your admin account.</template>
            <template v-else-if="step === 'totp'">Enter your authenticator code to continue.</template>
            <template v-else>Scan the QR with your authenticator app.</template>
          </p>
        </div>
        <!-- Step 1: login -->
        <form v-if="step === 'login'" class="flex flex-col gap-4" @submit.prevent="submitLogin">
          <div class="flex flex-col gap-2">
            <Label for="email">Email</Label>
            <Input id="email" v-model="email" type="email" placeholder="you@example.com" />
          </div>
          <div class="flex flex-col gap-2">
            <Label for="password">Password</Label>
            <div class="relative">
              <Input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                class="pr-10"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 flex items-center px-3 text-muted-foreground hover:text-foreground"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="showPassword" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
                  <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68" />
                  <path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
                  <line x1="2" x2="22" y1="2" y2="22" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
              </button>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox id="remember" :checked="rememberMe" @update:checked="rememberMe = $event" />
            <Label for="remember" class="text-sm font-normal">Remember me</Label>
          </div>
          <Button type="submit" :disabled="loading" class="w-full">
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </Button>
        </form>

        <!-- Step 2: totp -->
        <form v-else-if="step === 'totp'" class="flex flex-col gap-6" @submit.prevent="submitTotp">
          <div class="flex justify-center gap-2">
            <input
              v-for="(_, i) in 6"
              :key="i"
              :ref="(el: any) => { if (el) otpInputs[i] = el }"
              v-model="otp[i]"
              type="text"
              inputmode="numeric"
              autocomplete="one-time-code"
              maxlength="1"
              class="h-12 w-11 rounded-md border border-input bg-background text-center font-mono text-lg tabular-nums text-foreground transition-colors focus:border-ring focus:outline-none focus:ring-2 focus:ring-ring sm:w-12"
              @input="onOtpInput($event, i)"
              @keydown="onOtpKeydown($event, i)"
              @paste="onOtpPaste($event, i)"
            />
          </div>
          <Button type="submit" :disabled="loading" class="w-full">
            {{ loading ? 'Verifying…' : 'Verify' }}
          </Button>
        </form>

        <!-- Step 3: totp-setup -->
        <div v-else class="flex flex-col gap-4">
          <p class="text-sm text-muted-foreground">
            Scan the QR code with your authenticator app, then enter the 6-digit code.
          </p>
          <div class="flex justify-center rounded-lg border bg-white p-2">
            <img :src="totpQrUrl" alt="TOTP QR code" class="h-44 w-44" />
          </div>
          <div class="rounded-lg border p-2 text-center font-mono text-xs break-all">
            {{ totpSecret }}
          </div>
          <form class="flex flex-col gap-2" @submit.prevent="submitTotpSetup">
            <Input v-model="code" placeholder="123456" inputmode="numeric" />
            <Button type="submit" :disabled="loading">
              {{ loading ? 'Verifying…' : 'Confirm setup' }}
            </Button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
