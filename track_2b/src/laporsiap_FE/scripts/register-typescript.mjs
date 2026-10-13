// Test-only transpilation of actual production TypeScript; no app bundling.
import { registerHooks } from 'node:module'
import fs from 'node:fs'
import ts from 'typescript'
registerHooks({
 resolve(specifier,context,next){try{return next(specifier,context)}catch(error){if(specifier.startsWith('.'))for(const extension of ['.ts','.tsx'])try{return next(specifier+extension,context)}catch{}throw error}},
 load(url,context,next){if(/\.tsx?$/.test(url)){const source=fs.readFileSync(new URL(url),'utf8');return {format:'module',source:ts.transpileModule(source,{compilerOptions:{module:ts.ModuleKind.ESNext,target:ts.ScriptTarget.ES2022,jsx:ts.JsxEmit.ReactJSX,verbatimModuleSyntax:true}}).outputText,shortCircuit:true}}return next(url,context)},
})
