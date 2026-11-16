// Only the newest requested resource may be delivered to its caller. Repeated
// requests for the in-flight key share the same operation; completed data is not cached.
export function createOnDemandLoader<T>(read: (key: string, signal: AbortSignal) => Promise<T>) {
  let pending: { key: string; controller: AbortController; promise: Promise<T | undefined> } | undefined

  const cancel = () => {
    pending?.controller.abort()
    pending = undefined
  }

  const load = (key: string): Promise<T | undefined> => {
    if (pending?.key === key) return pending.promise
    cancel()
    const controller = new AbortController()
    const promise = Promise.resolve().then(() => {
      controller.signal.throwIfAborted()
      return read(key, controller.signal)
    }).then(
      (value) => controller.signal.aborted ? undefined : value,
      (error: unknown) => {
        if (!controller.signal.aborted) throw error
        return undefined
      }
    ).finally(() => {
      if (pending?.controller === controller) pending = undefined
    })
    pending = { key, controller, promise }
    return promise
  }

  return { load, cancel }
}
