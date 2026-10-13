import React from 'react'
import ReactDOM from 'react-dom/client'
import { Provider } from 'react-redux'
import App from './presentation/App'
import { LanguageProvider } from './presentation/i18n'
import { WorkflowProvider } from './presentation/services'
import { store, workflow, preferences } from './app/bootstrap'
import './style.css'
ReactDOM.createRoot(document.getElementById('root')!).render(<React.StrictMode><Provider store={store}><WorkflowProvider workflow={workflow}><LanguageProvider preferences={preferences}><App /></LanguageProvider></WorkflowProvider></Provider></React.StrictMode>)
