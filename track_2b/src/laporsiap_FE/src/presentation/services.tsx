import { createContext, useContext } from 'react'
import type { ReactNode } from 'react'
import type { CaseWorkflow } from '../application/caseWorkflow'
const Context = createContext<CaseWorkflow | null>(null)
export function WorkflowProvider({ workflow, children }: { workflow: CaseWorkflow; children: ReactNode }) { return <Context.Provider value={workflow}>{children}</Context.Provider> }
export function useWorkflow() { const workflow = useContext(Context); if (!workflow) throw new Error('WorkflowProvider is required.'); return workflow }
