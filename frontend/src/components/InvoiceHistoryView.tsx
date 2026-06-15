import { useState, useEffect } from 'react';
import { History, Building2, Calendar, FileText, Euro } from 'lucide-react';
import type { Client } from './ClientsView';

interface Invoice {
  ID: number;
  InvoiceNumber: string;
  IssueDate: string;
  Total: number;
  Status: string;
}

export default function InvoiceHistoryView() {
  const [clients, setClients] = useState<Client[]>([]);
  const [selectedClient, setSelectedClient] = useState<string>('');
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Fetch clients for the dropdown
    const fetchClients = async () => {
      try {
        const response = await fetch('/api/clients');
        if (response.ok) {
          const data = await response.json();
          setClients(data);
        }
      } catch (error) {
        console.error(error);
      }
    };
    fetchClients();
  }, []);

  useEffect(() => {
    if (!selectedClient) {
      setInvoices([]);
      return;
    }

    const fetchInvoices = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`/api/clients/${selectedClient}/invoices`);
        if (response.ok) {
          const data = await response.json();
          setInvoices(data);
        }
      } catch (error) {
        console.error(error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchInvoices();
  }, [selectedClient]);

  const formatDate = (dateString: string) => {
    const d = new Date(dateString);
    return new Intl.DateTimeFormat('es-ES', { 
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit'
    }).format(d);
  };

  return (
    <div className="w-full bg-slate-900 border border-slate-800/60 rounded-2xl p-8 shadow-2xl">
      <div className="flex items-center space-x-3 mb-8">
        <div className="bg-purple-500/10 p-2.5 rounded-xl border border-purple-500/20">
          <History className="w-6 h-6 text-purple-400" />
        </div>
        <h2 className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-pink-400">
          Historial de Facturas
        </h2>
      </div>

      <div className="mb-8 max-w-md">
        <label className="flex items-center text-sm font-medium text-slate-400 mb-2">
          <Building2 className="w-4 h-4 mr-2 text-slate-500" />
          Selecciona un cliente para ver su historial
        </label>
        <div className="relative group">
          <select
            value={selectedClient}
            onChange={(e) => setSelectedClient(e.target.value)}
            className="w-full bg-slate-950/50 border border-slate-800 rounded-xl px-4 py-3.5 appearance-none focus:outline-none focus:ring-2 focus:ring-purple-500/50 focus:border-purple-500 transition-all text-slate-200"
          >
            <option value="" disabled>Selecciona un cliente...</option>
            {clients.map(client => (
              <option key={client.ID} value={client.ID}>{client.Name} ({client.TaxID})</option>
            ))}
          </select>
          <div className="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-slate-500">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>
      </div>

      {!selectedClient ? (
        <div className="text-center py-16 border-2 border-dashed border-slate-800 rounded-xl bg-slate-900/50">
          <FileText className="w-12 h-12 text-slate-700 mx-auto mb-3" />
          <p className="text-slate-400 font-medium">Selecciona un cliente arriba</p>
          <p className="text-slate-500 text-sm mt-1">Para visualizar las facturas emitidas a su nombre.</p>
        </div>
      ) : isLoading ? (
        <div className="text-center py-10 text-slate-500">Cargando facturas...</div>
      ) : invoices.length === 0 ? (
        <div className="text-center py-16 border-2 border-dashed border-slate-800 rounded-xl bg-slate-900/50">
          <FileText className="w-12 h-12 text-slate-700 mx-auto mb-3" />
          <p className="text-slate-400 font-medium">Este cliente no tiene facturas</p>
          <p className="text-slate-500 text-sm mt-1">Las facturas emitidas aparecerán aquí.</p>
        </div>
      ) : (
        <div className="overflow-x-auto rounded-xl border border-slate-800">
          <table className="w-full text-left text-sm text-slate-400">
            <thead className="text-xs text-slate-500 uppercase bg-slate-950/80 border-b border-slate-800">
              <tr>
                <th className="px-6 py-4">
                  <div className="flex items-center"><FileText className="w-4 h-4 mr-2"/> Nº Factura</div>
                </th>
                <th className="px-6 py-4">
                  <div className="flex items-center"><Calendar className="w-4 h-4 mr-2"/> Fecha Emisión</div>
                </th>
                <th className="px-6 py-4">Estado</th>
                <th className="px-6 py-4 text-right">
                  <div className="flex items-center justify-end"><Euro className="w-4 h-4 mr-2"/> Total</div>
                </th>
              </tr>
            </thead>
            <tbody>
              {invoices.map((inv) => (
                <tr key={inv.ID} className="border-b border-slate-800/50 hover:bg-slate-800/30 transition-colors last:border-0">
                  <td className="px-6 py-4 font-mono font-medium text-slate-200">{inv.InvoiceNumber}</td>
                  <td className="px-6 py-4">{formatDate(inv.IssueDate)}</td>
                  <td className="px-6 py-4">
                    <span className="px-2.5 py-1 rounded-full text-xs font-medium bg-slate-800 text-slate-300">
                      {inv.Status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right font-medium text-emerald-400">
                    {inv.Total.toFixed(2)} €
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
