import assert from 'node:assert/strict'
import test from 'node:test'
import { createOnDemandLoader } from '../src/onDemand.ts'

const deferred = <T>() => {
  let resolve!: (value: T) => void
  let reject!: (reason: Error) => void
  const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no })
  return { promise, resolve, reject }
}

test('loads only on demand, deduplicates pending requests, and does not poll', async (t) => {
  t.mock.timers.enable({ apis: ['setTimeout', 'setInterval'] })
  let calls = 0
  const result = deferred<string>()
  const loader = createOnDemandLoader(async () => { calls++; return result.promise })
  t.mock.timers.tick(60000)
  assert.equal(calls, 0)
  const first = loader.load('a')
  assert.equal(loader.load('a'), first)
  await Promise.resolve()
  assert.equal(calls, 1)
  result.resolve('project a')
  assert.equal(await first, 'project a')
  t.mock.timers.tick(60000)
  assert.equal(calls, 1)
  await loader.load('a')
  assert.equal(calls, 2, 'explicit refresh requests fresh contents')
})

test('rapid selection aborts old reads and discards late results even if transport ignores abort', async () => {
  const a = deferred<string>()
  const b = deferred<string>()
  const signals: AbortSignal[] = []
  const loader = createOnDemandLoader((key, signal) => {
    signals.push(signal)
    return key === 'a' ? a.promise : b.promise
  })
  const old = loader.load('a')
  await Promise.resolve()
  const latest = loader.load('b')
  await Promise.resolve()
  assert.equal(signals[0].aborted, true)
  b.resolve('b')
  assert.equal(await latest, 'b')
  a.resolve('a')
  assert.equal(await old, undefined)
})

test('errors are retryable; canceled failures do not reach another project', async () => {
  const result = deferred<string>()
  const loader = createOnDemandLoader(() => result.promise)
  const old = loader.load('a')
  await Promise.resolve()
  loader.cancel()
  result.reject(new Error('old failure'))
  assert.equal(await old, undefined)
  let failed = true
  const retry = createOnDemandLoader(async () => {
    if (failed) throw new Error('Project not found')
    return 'restored'
  })
  await assert.rejects(retry.load('a'), /Project not found/)
  failed = false
  assert.equal(await retry.load('a'), 'restored')
})

