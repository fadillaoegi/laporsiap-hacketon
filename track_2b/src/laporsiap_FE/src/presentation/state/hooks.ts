import { useCallback } from 'react'
import { useDispatch, useSelector, useStore } from 'react-redux'
import type { AppDispatch, AppStore, RootState } from './store'
import type { WorkspaceState } from './workspaceSlice'
import { workspaceActions } from './workspaceSlice'
export const useAppDispatch = useDispatch.withTypes<AppDispatch>()
export const useAppSelector = useSelector.withTypes<RootState>()
export const useAppStore = useStore.withTypes<AppStore>()
type Update<T> = T | ((old: T) => T)
// Functional updates execute in a thunk, so action payloads remain serializable.
export function useWorkspaceField<K extends keyof WorkspaceState>(key: K) {
 const value = useAppSelector(state => state.workspace[key])
 const dispatch = useAppDispatch()
 const set = useCallback((update: Update<WorkspaceState[K]>) => {
  dispatch((dispatch, getState) => {
   const value = typeof update === 'function' ? (update as (old: WorkspaceState[K]) => WorkspaceState[K])(getState().workspace[key]) : update
   dispatch(workspaceActions.patch({ [key]: value }))
  })
 }, [dispatch, key])
 return [value, set] as const
}
