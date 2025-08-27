import { Button } from "@/components/ui/button"
import { Database, Shield, Zap, Terminal, ArrowRight, CheckCircle, Sparkles, Code, Globe } from "lucide-react"

function App() {
  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden">
      {/* Animated background gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/20 via-blue-900/20 to-teal-900/20"></div>
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(120,119,198,0.1),transparent_50%)]"></div>
      
      {/* Grid pattern overlay */}
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,.02)_1px,transparent_1px)] bg-[size:100px_100px]"></div>

      {/* Header */}
      <header className="relative z-50 border-b border-white/10 bg-black/50 backdrop-blur-xl">
        <div className="container mx-auto flex h-16 items-center justify-between px-6">
          <div className="flex items-center space-x-3">
            <div className="relative">
              <Database className="h-8 w-8 text-blue-400" />
              <div className="absolute -top-1 -right-1 h-3 w-3 bg-gradient-to-r from-blue-400 to-purple-400 rounded-full animate-pulse"></div>
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              db.xyz
            </span>
          </div>
          <nav className="hidden md:flex items-center space-x-8">
            <a href="#features" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
              Features
            </a>
            <a href="#pricing" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
              Pricing
            </a>
            <a href="#docs" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
              Docs
            </a>
            <Button variant="ghost" size="sm" className="text-gray-300 hover:text-white hover:bg-white/5">
              Sign In
            </Button>
            <Button size="sm" className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 border-0">
              Get Started
            </Button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <main className="relative z-10">
        <section className="container mx-auto px-6 pt-24 pb-16">
          <div className="text-center max-w-5xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-white/10 bg-white/5 backdrop-blur-sm mb-8">
              <Sparkles className="h-4 w-4 text-yellow-400" />
              <span className="text-sm text-gray-300">Now in Private Beta</span>
              <ArrowRight className="h-4 w-4 text-gray-400" />
            </div>
            
            <h1 className="text-6xl md:text-8xl font-medium mb-8 leading-tight">
              <span className="bg-gradient-to-r from-white via-blue-100 to-purple-100 bg-clip-text text-transparent">
                Postgres
              </span>
              <br />
              <span className="bg-gradient-to-r from-blue-400 via-purple-400 to-teal-400 bg-clip-text text-transparent">
                Reimagined
              </span>
            </h1>
            
            <p className="text-xl md:text-2xl text-gray-300 mb-12 max-w-3xl mx-auto leading-relaxed font-light">
              PostgreSQL databases deployed in seconds.
              <br />
              <span className="text-gray-400">Hardened containers. Zero configuration.</span>
            </p>
            
            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-16">
              <Button size="lg" className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 border-0 text-lg px-8 py-4 h-auto">
                Start Building
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button variant="outline" size="lg" className="border-white/20 bg-white/5 hover:bg-white/10 text-lg px-8 py-4 h-auto backdrop-blur-sm">
                <Terminal className="mr-2 h-5 w-5" />
                View CLI
              </Button>
            </div>

            {/* Code snippet preview */}
            <div className="max-w-2xl mx-auto">
              <div className="bg-black/40 backdrop-blur-sm border border-white/10 rounded-2xl p-6 text-left">
                <div className="flex items-center gap-2 mb-4">
                  <div className="flex gap-2">
                    <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <span className="text-sm text-gray-400 ml-4">Terminal</span>
                </div>
                <div className="font-mono text-sm space-y-2">
                  <div className="text-gray-400">$ npm install -g @db.xyz/cli</div>
                  <div className="text-gray-400">$ dbx login</div>
                  <div className="text-blue-400">$ dbx instance create --name prod-db --plan pro</div>
                  <div className="text-green-400">✓ Database created in 12s</div>
                  <div className="text-gray-300">🔗 postgresql://user:pass@pg-abc123.cust.db.xyz:5432/prod</div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Social proof / logos section */}
        <section className="container mx-auto px-6 py-16 border-t border-white/5">
          <div className="text-center mb-12">
            <p className="text-gray-400 text-sm uppercase tracking-wider">Trusted by developers worldwide</p>
          </div>
          <div className="flex justify-center items-center gap-12 opacity-40">
            {['Acme Corp', 'BuildCo', 'DevStudio', 'TechFlow', 'DataLab'].map((company) => (
              <div key={company} className="text-gray-500 font-medium text-lg">
                {company}
              </div>
            ))}
          </div>
        </section>

        {/* Features */}
        <section id="features" className="container mx-auto px-6 py-24">
          <div className="text-center mb-20">
            <h2 className="text-4xl md:text-5xl font-normal mb-6 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Built for Performance
            </h2>
            <p className="text-xl text-gray-400 max-w-3xl mx-auto font-light">
              Every database instance runs in an isolated, hardened container with built-in security and monitoring.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 mb-20">
            {[
              {
                icon: Shield,
                title: "Zero-Trust Security",
                description: "Unprivileged LXC containers with network isolation, encrypted connections, and automatic security updates.",
                gradient: "from-blue-400 to-cyan-400"
              },
              {
                icon: Zap,
                title: "Deploy in Seconds",
                description: "From API call to running database in under 30 seconds. Scale up or down instantly without downtime.",
                gradient: "from-purple-400 to-pink-400"
              },
              {
                icon: Code,
                title: "Developer First",
                description: "CLI tools, REST API, and web console. Integrate with your existing workflow and CI/CD pipelines.",
                gradient: "from-green-400 to-teal-400"
              }
            ].map((feature, index) => (
              <div key={index} className="group p-8 rounded-2xl border border-white/10 bg-gradient-to-br from-white/5 to-white/[0.02] hover:from-white/10 hover:to-white/5 transition-all duration-500">
                <div className={`inline-flex p-3 rounded-xl bg-gradient-to-r ${feature.gradient} mb-6`}>
                  <feature.icon className="h-6 w-6 text-white" />
                </div>
                <h3 className="text-xl font-medium mb-4 text-white group-hover:text-gray-100 transition-colors">
                  {feature.title}
                </h3>
                <p className="text-gray-400 leading-relaxed">
                  {feature.description}
                </p>
              </div>
            ))}
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" className="container mx-auto px-6 py-24">
          <div className="text-center mb-20">
            <h2 className="text-4xl md:text-5xl font-normal mb-6 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Pricing
            </h2>
            <p className="text-xl text-gray-400 font-light">
              Pay only for what you use. No hidden fees.
            </p>
          </div>

          <div className="grid md:grid-cols-4 gap-6 max-w-7xl mx-auto">
            {[
              { 
                name: 'Nano', 
                price: '$9', 
                cpu: '1 vCPU', 
                ram: '1 GB', 
                disk: '20 GB', 
                connections: '100',
                popular: false
              },
              { 
                name: 'Lite', 
                price: '$29', 
                cpu: '2 vCPU', 
                ram: '4 GB', 
                disk: '100 GB', 
                connections: '200',
                popular: false
              },
              { 
                name: 'Pro', 
                price: '$89', 
                cpu: '4 vCPU', 
                ram: '8 GB', 
                disk: '200 GB', 
                connections: '400',
                popular: true
              },
              { 
                name: 'Pro Heavy', 
                price: '$189', 
                cpu: '8 vCPU', 
                ram: '16 GB', 
                disk: '400 GB', 
                connections: '800',
                popular: false
              },
            ].map((plan, index) => (
              <div key={index} className="p-6 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 transition-all duration-300">
                <div>
                  <h3 className="text-lg font-medium mb-2 text-white">{plan.name}</h3>
                  <div className="mb-6">
                    <span className="text-3xl font-normal text-white">{plan.price}</span>
                    <span className="text-gray-400 text-sm">/month</span>
                  </div>
                  
                  <div className="space-y-2 mb-6 text-sm">
                    <div className="text-gray-300">{plan.cpu}</div>
                    <div className="text-gray-300">{plan.ram} RAM</div>
                    <div className="text-gray-300">{plan.disk} SSD</div>
                    <div className="text-gray-300">{plan.connections} connections</div>
                  </div>
                  
                  <Button className="w-full bg-white/10 hover:bg-white/20 border-0 text-sm">
                    Get Started
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </section>

      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-white/10 bg-black/50 backdrop-blur-xl">
        <div className="container mx-auto px-6 py-12">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="flex items-center space-x-3 mb-4">
                <Database className="h-6 w-6 text-blue-400" />
                <span className="text-lg font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                  db.xyz
                </span>
              </div>
              <p className="text-gray-400 text-sm leading-relaxed font-light">
                Postgres-as-a-Service on hardened LXC containers.
              </p>
            </div>
            
            <div>
              <h4 className="font-medium text-white mb-3">Product</h4>
              <div className="space-y-2">
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Features</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Pricing</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Security</a>
              </div>
            </div>
            
            <div>
              <h4 className="font-medium text-white mb-3">Resources</h4>
              <div className="space-y-2">
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Documentation</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">API Reference</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">CLI Guide</a>
              </div>
            </div>
            
            <div>
              <h4 className="font-medium text-white mb-3">Company</h4>
              <div className="space-y-2">
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">About</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Status</a>
                <a href="#" className="block text-gray-400 hover:text-white text-sm transition-colors">Support</a>
              </div>
            </div>
          </div>
          
          <div className="border-t border-white/10 pt-8 flex flex-col md:flex-row justify-between items-center">
            <p className="text-gray-400 text-sm">
              © 2024 db.xyz. All rights reserved.
            </p>
            <div className="flex gap-6 mt-4 md:mt-0">
              <a href="#" className="text-gray-400 hover:text-white text-sm transition-colors">Privacy</a>
              <a href="#" className="text-gray-400 hover:text-white text-sm transition-colors">Terms</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App