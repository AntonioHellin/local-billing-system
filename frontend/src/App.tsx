import { useState, useEffect } from 'react';
import { Toaster, toast } from 'react-hot-toast';
import { FileText, Send, Building2, AlignLeft, Euro, Receipt, Users, History } from 'lucide-react';
import ClientsView from './components/ClientsView';
import InvoiceHistoryView from './components/InvoiceHistoryView';

// Type definitions based on Go models
import type { Client } from './components/ClientsView';



interface InvoicePayload {
  client_id: number;
  base_imponible: number;
}

function App() {
  const [activeTab, setActiveTab] = useState<'invoice' | 'clients' | 'history'>('invoice');

  // Form states for New Invoice
  const [clients, setClients] = useState<Client[]>([]);
  const [selectedClient, setSelectedClient] = useState<string>('');
  const [concept, setConcept] = useState<string>('');
  const [baseAmount, setBaseAmount] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // VAT is 21%
  const TAX_RATE = 0.21;
  
  const subtotal = parseFloat(baseAmount) || 0;
  const taxTotal = subtotal * TAX_RATE;
  const total = subtotal + taxTotal;

  useEffect(() => {
    // Only fetch clients when on invoice tab to populate dropdown
    if (activeTab === 'invoice') {
      const fetchClients = async () => {
        try {
          const response = await fetch('/api/clients');
          if (response.ok) {
            const data = await response.json();
            setClients(data);
          }
        } catch (error) {
          console.error('Error fetching clients', error);
        }
      };

      fetchClients();
    }
  }, [activeTab]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!selectedClient || !concept || subtotal <= 0) {
      toast.error('Por favor, completa todos los campos correctamente.');
      return;
    }

    setIsSubmitting(true);
    
    const payload: InvoicePayload = {
      client_id: parseInt(selectedClient, 10),
      base_imponible: subtotal,
    };

    try {
      // Simulate API call delay for effect
      await new Promise(resolve => setTimeout(resolve, 800));
      
      const response = await fetch('/api/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        throw new Error('Error al generar la factura');
      }

      toast.success('¡Factura emitida con éxito!');
      
      // Reset form
      setSelectedClient('');
      setConcept('');
      setBaseAmount('');
    } catch (error) {
      toast.error('Hubo un problema al emitir la factura.');
      console.error(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 flex flex-col font-sans">
      <Toaster position="top-right" toastOptions={{ 
        className: 'bg-slate-800 text-slate-100 border border-slate-700',
        style: {
          background: '#1e293b',
          color: '#f1f5f9',
          border: '1px solid #334155'
        }
      }} />

      {/* Top Navigation */}
      <nav className="bg-slate-900 border-b border-slate-800 px-6 py-4 sticky top-0 z-50 shadow-md">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="bg-gradient-to-br from-indigo-500 to-purple-600 p-2 rounded-xl text-white">
              <Receipt className="w-5 h-5" />
            </div>
            <span className="text-xl font-bold text-slate-100 tracking-tight">FacturaLocal</span>
          </div>
          
          <div className="flex space-x-2">
            <button
              onClick={() => setActiveTab('invoice')}
              className={`flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                activeTab === 'invoice' 
                  ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' 
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
              }`}
            >
              <FileText className="w-4 h-4 mr-2" /> Nueva Factura
            </button>
            <button
              onClick={() => setActiveTab('clients')}
              className={`flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                activeTab === 'clients' 
                  ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' 
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
              }`}
            >
              <Users className="w-4 h-4 mr-2" /> Clientes
            </button>
            <button
              onClick={() => setActiveTab('history')}
              className={`flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                activeTab === 'history' 
                  ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' 
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
              }`}
            >
              <History className="w-4 h-4 mr-2" /> Historial
            </button>
          </div>
        </div>
      </nav>

      {/* Main Content Area */}
      <main className="flex-1 max-w-7xl w-full mx-auto p-4 lg:p-8 flex items-start justify-center">
        {activeTab === 'invoice' && (
          <div className="w-full max-w-5xl bg-slate-900 border border-slate-800/60 rounded-2xl shadow-2xl overflow-hidden flex flex-col lg:flex-row backdrop-blur-sm">
        
        {/* Left Form Column */}
        <div className="w-full lg:w-3/5 p-8 lg:p-12">
          <div className="flex items-center space-x-3 mb-8">
            <div className="bg-indigo-500/10 p-2.5 rounded-xl border border-indigo-500/20">
              <FileText className="w-6 h-6 text-indigo-400" />
            </div>
            <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-400 to-purple-400">
              Nueva Factura
            </h1>
          </div>

          <form onSubmit={handleSubmit} className="space-y-8">
            {/* Client Selection */}
            <div className="space-y-3">
              <label htmlFor="client" className="flex items-center text-sm font-medium text-slate-400">
                <Building2 className="w-4 h-4 mr-2 text-slate-500" />
                Cliente
              </label>
              <div className="relative group">
                <select
                  id="client"
                  value={selectedClient}
                  onChange={(e) => setSelectedClient(e.target.value)}
                  className="w-full bg-slate-950/50 border border-slate-800 rounded-xl px-4 py-3.5 appearance-none focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all duration-200 text-slate-200"
                >
                  <option value="" disabled>Selecciona un cliente...</option>
                  {clients.map(client => (
                    <option key={client.ID} value={client.ID}>{client.Name}</option>
                  ))}
                </select>
                <div className="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-slate-500 group-hover:text-indigo-400 transition-colors">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path></svg>
                </div>
              </div>
            </div>

            {/* Concept Description */}
            <div className="space-y-3">
              <label htmlFor="concept" className="flex items-center text-sm font-medium text-slate-400">
                <AlignLeft className="w-4 h-4 mr-2 text-slate-500" />
                Concepto de Facturación
              </label>
              <textarea
                id="concept"
                value={concept}
                onChange={(e) => setConcept(e.target.value)}
                placeholder="Ej. Servicios de consultoría técnica para el mes de mayo..."
                rows={4}
                className="w-full bg-slate-950/50 border border-slate-800 rounded-xl px-4 py-3.5 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all duration-200 text-slate-200 resize-none placeholder-slate-600"
              />
            </div>

            {/* Base Amount */}
            <div className="space-y-3">
              <label htmlFor="baseAmount" className="flex items-center text-sm font-medium text-slate-400">
                <Euro className="w-4 h-4 mr-2 text-slate-500" />
                Base Imponible
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
                  <span className="text-slate-500 font-medium">€</span>
                </div>
                <input
                  id="baseAmount"
                  type="number"
                  step="0.01"
                  min="0"
                  value={baseAmount}
                  onChange={(e) => setBaseAmount(e.target.value)}
                  placeholder="0.00"
                  className="w-full bg-slate-950/50 border border-slate-800 rounded-xl pl-10 pr-4 py-3.5 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all duration-200 text-slate-200 placeholder-slate-600 font-medium"
                />
              </div>
            </div>
          </form>
        </div>

        {/* Right Sidebar Summary */}
        <div className="w-full lg:w-2/5 bg-slate-800/30 p-8 lg:p-12 border-t lg:border-t-0 lg:border-l border-slate-800 flex flex-col justify-between">
          
          <div>
            <div className="flex items-center space-x-2 mb-8">
              <Receipt className="w-5 h-5 text-indigo-400" />
              <h2 className="text-xl font-semibold text-slate-200">Resumen</h2>
            </div>
            
            <div className="space-y-6">
              <div className="flex justify-between items-center py-3 border-b border-slate-800/60">
                <span className="text-slate-400">Base Imponible</span>
                <span className="font-medium text-slate-200">{subtotal.toFixed(2)} €</span>
              </div>
              
              <div className="flex justify-between items-center py-3 border-b border-slate-800/60">
                <span className="text-slate-400">IVA (21%)</span>
                <span className="font-medium text-slate-200">{taxTotal.toFixed(2)} €</span>
              </div>
              
              <div className="flex justify-between items-center py-4">
                <span className="text-lg font-medium text-slate-300">Total a Pagar</span>
                <span className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 to-cyan-400">
                  {total.toFixed(2)} €
                </span>
              </div>
            </div>
          </div>

          <div className="mt-12 pt-8 border-t border-slate-800/60">
            <button
              onClick={handleSubmit}
              disabled={isSubmitting || !selectedClient || !concept || subtotal <= 0}
              className="w-full flex items-center justify-center py-4 px-6 rounded-xl text-white font-medium bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-400 hover:to-purple-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40 group transform hover:-translate-y-0.5 active:translate-y-0"
            >
              {isSubmitting ? (
                <div className="flex items-center">
                  <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Procesando...
                </div>
              ) : (
                <div className="flex items-center">
                  <Send className="w-5 h-5 mr-2 group-hover:translate-x-1 transition-transform" />
                  Emitir Factura
                </div>
              )}
            </button>
          </div>

        </div>
      </div>
        )}

        {activeTab === 'clients' && <div className="w-full"><ClientsView /></div>}
        {activeTab === 'history' && <div className="w-full max-w-5xl mx-auto"><InvoiceHistoryView /></div>}
      </main>
    </div>
  );
}

export default App;
