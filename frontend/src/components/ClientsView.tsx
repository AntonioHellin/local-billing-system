import { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { Users, UserPlus, FileSignature, Mail, Phone, MapPin } from 'lucide-react';

export interface Client {
  ID: number;
  Name: string;
  TaxID: string;
  Email?: string;
  Phone?: string;
  Address?: string;
}

export default function ClientsView() {
  const [clients, setClients] = useState<Client[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Form states
  const [name, setName] = useState('');
  const [taxId, setTaxId] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [address, setAddress] = useState('');

  const fetchClients = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('/api/clients');
      if (response.ok) {
        const data = await response.json();
        setClients(data);
      }
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchClients();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !taxId) {
      toast.error('Nombre y NIF/CIF son requeridos.');
      return;
    }

    setIsSubmitting(true);
    try {
      const response = await fetch('/api/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name,
          tax_id: taxId,
          email,
          phone,
          address
        })
      });

      if (!response.ok) throw new Error('Error al crear el cliente');

      toast.success('Cliente creado con éxito');
      setName('');
      setTaxId('');
      setEmail('');
      setPhone('');
      setAddress('');
      
      // Refresh list
      fetchClients();
    } catch (error) {
      toast.error('Hubo un problema al crear el cliente.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col lg:flex-row gap-8">
      {/* Create Client Form */}
      <div className="w-full lg:w-1/3 bg-slate-900 border border-slate-800/60 rounded-2xl p-6 shadow-2xl">
        <div className="flex items-center space-x-3 mb-6">
          <div className="bg-emerald-500/10 p-2.5 rounded-xl border border-emerald-500/20">
            <UserPlus className="w-5 h-5 text-emerald-400" />
          </div>
          <h2 className="text-xl font-bold text-slate-200">Nuevo Cliente</h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1">Nombre Completo *</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full bg-slate-950/50 border border-slate-800 rounded-lg px-3 py-2.5 text-sm focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500 transition-all text-slate-200"
              placeholder="Ej. Tech Corp S.L."
            />
          </div>

          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1">NIF / CIF *</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-500">
                <FileSignature className="w-4 h-4" />
              </div>
              <input
                type="text"
                value={taxId}
                onChange={(e) => setTaxId(e.target.value)}
                className="w-full bg-slate-950/50 border border-slate-800 rounded-lg pl-9 pr-3 py-2.5 text-sm focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500 transition-all text-slate-200 uppercase"
                placeholder="Ej. B12345678"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1">Email</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-500">
                <Mail className="w-4 h-4" />
              </div>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-slate-950/50 border border-slate-800 rounded-lg pl-9 pr-3 py-2.5 text-sm focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500 transition-all text-slate-200"
                placeholder="contacto@empresa.com"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1">Teléfono</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-500">
                <Phone className="w-4 h-4" />
              </div>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="w-full bg-slate-950/50 border border-slate-800 rounded-lg pl-9 pr-3 py-2.5 text-sm focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500 transition-all text-slate-200"
                placeholder="+34 600 000 000"
              />
            </div>
          </div>
          
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1">Dirección</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-500">
                <MapPin className="w-4 h-4" />
              </div>
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                className="w-full bg-slate-950/50 border border-slate-800 rounded-lg pl-9 pr-3 py-2.5 text-sm focus:ring-2 focus:ring-emerald-500/50 focus:border-emerald-500 transition-all text-slate-200"
                placeholder="Calle Mayor 1, Madrid"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={isSubmitting || !name || !taxId}
            className="w-full mt-4 py-2.5 rounded-lg text-white text-sm font-medium bg-emerald-600 hover:bg-emerald-500 focus:ring-2 focus:ring-emerald-500/50 transition-all disabled:opacity-50"
          >
            {isSubmitting ? 'Guardando...' : 'Registrar Cliente'}
          </button>
        </form>
      </div>

      {/* Clients List */}
      <div className="w-full lg:w-2/3 bg-slate-900 border border-slate-800/60 rounded-2xl p-6 shadow-2xl">
        <div className="flex items-center space-x-3 mb-6">
          <div className="bg-indigo-500/10 p-2.5 rounded-xl border border-indigo-500/20">
            <Users className="w-5 h-5 text-indigo-400" />
          </div>
          <h2 className="text-xl font-bold text-slate-200">Cartera de Clientes</h2>
        </div>

        {isLoading ? (
          <div className="text-center py-10 text-slate-500">Cargando clientes...</div>
        ) : clients.length === 0 ? (
          <div className="text-center py-12 border-2 border-dashed border-slate-800 rounded-xl">
            <Users className="w-12 h-12 text-slate-700 mx-auto mb-3" />
            <p className="text-slate-400 font-medium">No hay clientes registrados</p>
            <p className="text-slate-500 text-sm mt-1">Usa el formulario de la izquierda para añadir uno.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="text-xs text-slate-500 uppercase bg-slate-950/50">
                <tr>
                  <th className="px-4 py-3 rounded-tl-lg">Nombre</th>
                  <th className="px-4 py-3">NIF/CIF</th>
                  <th className="px-4 py-3">Email</th>
                  <th className="px-4 py-3 rounded-tr-lg">Teléfono</th>
                </tr>
              </thead>
              <tbody>
                {clients.map((client) => (
                  <tr key={client.ID} className="border-b border-slate-800/50 hover:bg-slate-800/20 transition-colors">
                    <td className="px-4 py-4 font-medium text-slate-200">{client.Name}</td>
                    <td className="px-4 py-4 font-mono text-xs">{client.TaxID}</td>
                    <td className="px-4 py-4">{client.Email || '-'}</td>
                    <td className="px-4 py-4">{client.Phone || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
