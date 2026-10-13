import fs from 'node:fs'
import path from 'node:path'
import ts from 'typescript'
const root=path.resolve(new URL('../src',import.meta.url).pathname)
const violations=[];let checked=0
const walk=folder=>{
 for(const entry of fs.readdirSync(folder,{withFileTypes:true})){
  const file=path.join(folder,entry.name);if(entry.isDirectory()){walk(file);continue};if(!/\.tsx?$/.test(file))continue
  const source=fs.readFileSync(file,'utf8'),layer=path.relative(root,file).split(path.sep)[0]
  const parsed=ts.createSourceFile(file,source,ts.ScriptTarget.Latest,true);checked++
  for(const statement of parsed.statements){
   if(!(ts.isImportDeclaration(statement)||ts.isExportDeclaration(statement))||!statement.moduleSpecifier||!ts.isStringLiteral(statement.moduleSpecifier))continue
   const dependency=statement.moduleSpecifier.text
   const target=dependency.startsWith('.')?path.relative(root,path.resolve(path.dirname(file),dependency)).split(path.sep)[0]:dependency
   const allowed={domain:['domain'],application:['application','domain'],infrastructure:['infrastructure','application','domain'],presentation:['presentation','application','domain','react','react-redux','@reduxjs/toolkit']}[layer]
   if(allowed&&!allowed.includes(target))violations.push(`${path.relative(root,file)} imports ${dependency}`)
  }
  if(['domain','application'].includes(layer)&&/\b(?:fetch\s*\(|(?:localStorage|sessionStorage|document|window)\s*\.)/.test(source.replace(/\/\/[^\n]*/g,'')))violations.push(`${path.relative(root,file)} directly accesses a browser adapter`)
 }
}
walk(root)
if(violations.length){console.error(violations.join('\n'));process.exit(1)}
console.log(`Architecture boundaries passed (${checked} TypeScript files).`)
