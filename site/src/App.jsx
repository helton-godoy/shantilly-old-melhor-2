import { useState } from 'react'
import { Button } from '@/components/ui/button.jsx'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card.jsx'
import { Badge } from '@/components/ui/badge.jsx'
import { Terminal, Code, Zap, Palette, Download, Github, BookOpen, Sparkles } from 'lucide-react'
import shantillyLogo from './assets/shantilly_v1.png'
import cupcakeTransparent from './assets/shantilly_cupcake_transparent.png'
import './App.css'

function App() {
  const [activeFeature, setActiveFeature] = useState(0)

  const features = [
    {
      icon: <Terminal className="w-6 h-6" />,
      title: "Interface de Terminal Moderna",
      description: "Crie TUIs (Terminal User Interfaces) interativas e visualmente atraentes com facilidade."
    },
    {
      icon: <Code className="w-6 h-6" />,
      title: "Configuração Declarativa",
      description: "Defina suas interfaces usando arquivos YAML ou JSON simples, sem código imperativo."
    },
    {
      icon: <Zap className="w-6 h-6" />,
      title: "Alta Performance",
      description: "Construído em Go com o Charm stack para máxima velocidade e eficiência."
    },
    {
      icon: <Palette className="w-6 h-6" />,
      title: "Estilização Avançada",
      description: "Personalize cores, layouts e animações com o poder do Lipgloss."
    }
  ]

  const components = [
    { name: "Formulários", description: "Campos de entrada interativos" },
    { name: "Menus", description: "Navegação intuitiva" },
    { name: "Layouts", description: "Organização flexível" },
    { name: "Abas", description: "Conteúdo organizado" },
    { name: "Tabelas", description: "Dados estruturados" },
    { name: "Barras de Progresso", description: "Feedback visual" }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      {/* Header */}
      <header className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-green-500/10 to-pink-500/10"></div>
        <div className="relative container mx-auto px-6 py-20">
          <div className="flex flex-col lg:flex-row items-center justify-between gap-12">
            <div className="flex-1 text-center lg:text-left">
              <div className="flex items-center justify-center lg:justify-start mb-6">
                <img src={cupcakeTransparent} alt="Shantilly Cupcake" className="h-16 w-auto mr-4" />
                <div className="text-4xl font-bold">
                  <span className="text-green-400">SH</span>
                  <span className="text-pink-400">antilly</span>
                </div>
              </div>
              <p className="text-xl lg:text-2xl text-slate-300 mb-8 font-medium">
                Cobertura gráfica para seu script
              </p>
              <p className="text-lg text-slate-400 mb-10 max-w-2xl">
                Transforme seus scripts de terminal em interfaces de usuário modernas e interativas.
                Shantilly oferece uma maneira simples e poderosa de criar TUIs profissionais usando
                apenas arquivos de configuração.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                <Button size="lg" className="bg-green-500 hover:bg-green-600 text-white px-8 py-3 text-lg" onClick={() => window.location.href = '/shantilly/docs/pt-br/'}>
                  <Download className="w-5 h-5 mr-2" />
                  Começar Agora
                </Button>
                <Button variant="outline" size="lg" className="border-slate-600 text-slate-300 hover:bg-slate-800 px-8 py-3 text-lg" onClick={() => window.open('https://github.com/helton-godoy/shantilly', '_blank')}>
                  <Github className="w-5 h-5 mr-2" />
                  Ver no GitHub
                </Button>
              </div>
            </div>
            <div className="flex-1 max-w-lg">
              {/* Terminal Simulation with CSS */}
              <div className="terminal-window bg-slate-900 rounded-lg overflow-hidden shadow-2xl">
                <div className="terminal-header bg-slate-800 px-4 py-2 flex items-center gap-2">
                  <div className="flex gap-2">
                    <div className="w-3 h-3 rounded-full bg-red-500"></div>
                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                  </div>
                  <div className="text-slate-400 text-sm ml-4">$ shantilly run config.yaml</div>
                </div>
                <div className="terminal-content p-6 font-mono text-sm">
                  <div className="text-green-400 mb-4">✓ Terminal Interativo</div>

                  {/* User Registration Form with CSS borders */}
                  <div className="terminal-form">
                    <div className="form-title text-green-400 mb-4 font-bold">User Registration</div>

                    <div className="form-container">
                      <div className="form-field mb-3">
                        <label className="text-slate-300 block mb-1">Username:</label>
                        <div className="input-field bg-slate-800 border border-green-400 rounded px-3 py-2 text-white">
                          admin
                        </div>
                      </div>

                      <div className="form-field mb-3">
                        <label className="text-slate-300 block mb-1">Password:</label>
                        <div className="input-field bg-slate-800 border border-slate-600 rounded px-3 py-2 text-white">
                          ••••••••
                        </div>
                      </div>

                      <div className="form-field mb-4">
                        <label className="text-slate-300 block mb-1">Email:</label>
                        <div className="input-field bg-slate-800 border border-slate-600 rounded px-3 py-2 text-slate-400">
                          user@example.com
                        </div>
                        <div className="text-red-400 text-xs mt-1">Email inválido.</div>
                      </div>

                      <div className="form-field mb-4">
                        <label className="text-slate-300 flex items-center">
                          <input type="checkbox" className="mr-2 accent-green-400" checked readOnly />
                          Accept Terms and Conditions
                        </label>
                      </div>

                      <div className="form-buttons flex gap-3">
                        <button className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded font-semibold">
                          Submit
                        </button>
                        <button className="bg-slate-600 hover:bg-slate-700 text-white px-4 py-2 rounded">
                          Cancel
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Features Section */}
      <section className="py-20 bg-slate-800/30">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">
              Por que escolher o Shantilly?
            </h2>
            <p className="text-xl text-slate-400 max-w-3xl mx-auto">
              Uma ferramenta moderna que revoluciona a criação de interfaces de terminal,
              combinando simplicidade com poder e performance.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <Card
                key={index}
                className="bg-slate-800/50 border-slate-700 hover:bg-slate-700/50 transition-all duration-300 cursor-pointer group"
                onClick={() => setActiveFeature(index)}
              >
                <CardHeader>
                  <div className="text-green-400 mb-4 group-hover:text-pink-400 transition-colors">
                    {feature.icon}
                  </div>
                  <CardTitle className="text-white text-lg">{feature.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-slate-400">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Components Section */}
      <section className="py-20">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">
              Componentes Disponíveis
            </h2>
            <p className="text-xl text-slate-400 max-w-3xl mx-auto">
              Uma biblioteca completa de componentes prontos para usar em suas interfaces de terminal.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {components.map((component, index) => (
              <Card key={index} className="bg-slate-800/50 border-slate-700 hover:border-green-500/50 transition-all duration-300">
                <CardHeader>
                  <CardTitle className="text-white flex items-center justify-between">
                    {component.name}
                    <Sparkles className="w-4 h-4 text-pink-400" />
                  </CardTitle>
                  <CardDescription className="text-slate-400">
                    {component.description}
                  </CardDescription>
                </CardHeader>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Technology Stack */}
      <section className="py-20 bg-slate-800/30">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">
              Tecnologias Modernas
            </h2>
            <p className="text-xl text-slate-400 max-w-3xl mx-auto">
              Construído sobre o robusto ecossistema Go com as melhores bibliotecas disponíveis.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { name: "Go", description: "Performance e confiabilidade" },
              { name: "Cobra", description: "Framework CLI moderno" },
              { name: "Bubble Tea", description: "Gerenciamento de estado TUI" },
              { name: "Lipgloss", description: "Estilização avançada" }
            ].map((tech, index) => (
              <div key={index} className="text-center">
                <Badge variant="outline" className="border-green-500 text-green-400 mb-3 px-4 py-2 text-lg">
                  {tech.name}
                </Badge>
                <p className="text-slate-400">{tech.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20">
        <div className="container mx-auto px-6 text-center">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-4xl lg:text-5xl font-bold text-white mb-6">
              Pronto para transformar seus scripts?
            </h2>
            <p className="text-xl text-slate-400 mb-10">
              Comece a criar interfaces de terminal profissionais em minutos, não horas.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" className="bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white px-8 py-4 text-lg" onClick={() => window.location.href = '/shantilly/docs/pt-br/'}>
                <Download className="w-5 h-5 mr-2" />
                Instalar Shantilly
              </Button>
              <Button variant="outline" size="lg" className="border-slate-600 text-slate-300 hover:bg-slate-800 px-8 py-4 text-lg" onClick={() => window.location.href = '/shantilly/docs/pt-br/'}>
                <BookOpen className="w-5 h-5 mr-2" />
                Ler Documentação
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 bg-slate-900 border-t border-slate-800">
        <div className="container mx-auto px-6">
          <div className="flex flex-col md:flex-row items-center justify-between">
            <div className="flex items-center mb-4 md:mb-0">
              <img src={cupcakeTransparent} alt="Shantilly" className="h-8 w-auto mr-3" />
              <span className="text-white font-semibold">
                <span className="text-green-400">SH</span>
                <span className="text-pink-400">antilly</span>
              </span>
            </div>
            <div className="flex items-center space-x-6">
              <a href="/shantilly/docs/pt-br/" className="text-slate-400 hover:text-white transition-colors">
                Documentação
              </a>
              <a href="https://github.com/helton-godoy/shantilly" target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-white transition-colors">
                GitHub
              </a>
              <a href="https://github.com/helton-godoy/shantilly/wiki" target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-white transition-colors">
                Wiki
              </a>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-slate-800 text-center">
            <p className="text-slate-500">
              © 2024 Shantilly. Cobertura gráfica para seu script.
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App
