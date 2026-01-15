import { useEffect } from 'react'

export function useResizeObserver<T extends HTMLElement>(
  ref: React.RefObject<T>,
  onSize: (width: number, height: number) => void
) {
  useEffect(() => {
    const el = ref.current
    if (!el) return
    let raf = 0
    const ro = new ResizeObserver(entries => {
      const rect = entries[0].contentRect
      cancelAnimationFrame(raf)
      raf = requestAnimationFrame(() => onSize(rect.width, rect.height))
    })
    ro.observe(el)
    return () => {
      cancelAnimationFrame(raf)
      ro.disconnect()
    }
  }, [ref, onSize])
}
