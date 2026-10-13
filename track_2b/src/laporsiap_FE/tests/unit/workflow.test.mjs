import test from 'node:test'
import assert from 'node:assert/strict'
import { CaseWorkflow } from '../../src/application/caseWorkflow.ts'
import { BrowserSessionVault } from '../../src/infrastructure/browser/sessionVault.ts'
import { createAppStore } from '../../src/presentation/state/store.ts'
import { workspaceActions } from '../../src/presentation/state/workspaceSlice.ts'
const caseData = (id='C1',revision=1)=>({id,title:'Fictional',synthetic:true,raw_chat:'I sent IDR 264,000.',parser:'plain',messages:[{id:'M1',text:'I sent IDR 264,000.',sender:'',sent_at:''}],transfers:[],draft:null,attachments:[],revision,reviewed_revision:null,reviewed_at:'',created_at:'',updated_at:'',expires_at:'',access_expires_at:''})
const summary={successful_known_total:null,successful_count:0,successful_unknown_amount:0,status_counts:{},warnings:[]}
test('session creation and rotation keep credentials outside Redux state/actions',async()=>{
 const storage=new Map();globalThis.sessionStorage={getItem:key=>storage.get(key)||null,setItem:(key,value)=>storage.set(key,value)}
 const vault=new BrowserSessionVault(),original='a'.repeat(48),rotated='b'.repeat(48),c=caseData()
 const repository={create:async()=>({case:c,summary,token:original}),rotate:async()=>({case:c,summary,token:rotated})}
 const workflow=new CaseWorkflow(repository,vault,{save(){}})
 const store=createAppStore();const input={title:c.title,raw_chat:c.raw_chat,parser:c.parser,synthetic:true,transfers:[]}
 const result=await workflow.save(null,input,'en');assert.equal(result.token,undefined);store.dispatch(workspaceActions.applyCase(result));store.dispatch(workspaceActions.patch({accessList:workflow.records()}))
 assert.equal(vault.get(c.id).token,original);assert.equal(JSON.stringify(store.getState()).includes(original),false)
 const next=await workflow.rotate(c,'en');store.dispatch(workspaceActions.applyCase(next));store.dispatch(workspaceActions.accessRotated())
 assert.equal(vault.get(c.id).token,rotated);assert.equal(JSON.stringify(store.getState()).includes(rotated),false);assert.equal(store.getState().workspace.accessEpoch,1)
})
test('failed case deletion preserves the session needed to retry manually',async()=>{
 const vault=new BrowserSessionVault();vault.remember({case:caseData(),summary,token:'c'.repeat(48)})
 const workflow=new CaseWorkflow({delete:async()=>{throw new Error('Offline')}},vault,{save(){}})
 await assert.rejects(workflow.delete(caseData(),'en'),/Offline/);assert.equal(vault.records().length,1)
})
test('new case and changed revision discard old privacy/review permissions',()=>{
 const store=createAppStore([{id:'C1',title:'Fictional'}]);store.dispatch(workspaceActions.applyCase({case:caseData(),summary}))
 store.dispatch(workspaceActions.patch({maskPrivate:false,privateTerms:'private',consent:true,aiPreview:{hash:'old'},reviewedSources:['M1'],reviewedTransfers:true,acceptedLimits:true,deleteConfirmation:'DELETE'}))
 store.dispatch(workspaceActions.applyCase({case:caseData('C1',2),summary}));assert.deepEqual(store.getState().workspace.reviewedSources,[]);assert.equal(store.getState().workspace.acceptedLimits,false)
 store.dispatch(workspaceActions.newCase());const s=store.getState().workspace;assert.equal(s.access,null);assert.equal(s.maskPrivate,true);assert.equal(s.privateTerms,'');assert.equal(s.consent,false);assert.equal(s.aiPreview,null);assert.equal(s.deleteConfirmation,'');assert.equal(s.accessList.length,1)
})
