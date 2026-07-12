<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'

const router = useRouter()

type Step = 'login' | 'totp' | 'totp-setup'
const step = ref<Step>('login')
const loading = ref(false)

const email = ref('')
const password = ref('')
const code = ref('')
const tempToken = ref('')
const totpSecret = ref('')
const totpQrUrl = ref('')

function storeTokenAndRedirect(token: string) {
  localStorage.setItem('access_token', token)
  router.push('/')
}

function storeToken(token: string) {
  localStorage.setItem('access_token', token)
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
    })
    if (data.requires_totp) {
      tempToken.value = data.temp_token
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

async function submitTotp() {
  if (!code.value) {
    toast.warning('Code is required')
    return
  }
  loading.value = true
  try {
    const { data } = await api.post('/auth/totp', {
      temp_token: tempToken.value,
      code: code.value,
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
  <div class="flex min-h-screen items-center justify-center bg-background p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle class="font-heading">okuru.id</CardTitle>
        <CardDescription>
          <template v-if="step === 'login'">Sign in to your admin account</template>
          <template v-else-if="step === 'totp'">Enter your authenticator code</template>
          <template v-else>Set up two-factor authentication</template>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <!-- Step 1: login -->
        <form v-if="step === 'login'" class="flex flex-col gap-4" @submit.prevent="submitLogin">
          <div class="flex flex-col gap-2">
            <Label for="email">Email</Label>
            <Input id="email" v-model="email" type="email" placeholder="you@example.com" />
          </div>
          <div class="flex flex-col gap-2">
            <Label for="password">Password</Label>
            <Input id="password" v-model="password" type="password" />
          </div>
          <Button type="submit" :disabled="loading" class="w-full">
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </Button>
        </form>

        <!-- Step 2: totp -->
        <form v-else-if="step === 'totp'" class="flex flex-col gap-4" @submit.prevent="submitTotp">
          <div class="flex flex-col gap-2">
            <Label for="code">Authentication code</Label>
            <Input
              id="code"
              v-model="code"
              placeholder="123456"
              autocomplete="one-time-code"
              inputmode="numeric"
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
      </CardContent>
    </Card>
  </div>
</template>
