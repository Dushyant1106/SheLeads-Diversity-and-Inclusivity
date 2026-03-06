import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/auth-context';
import ThemeToggle from '@/components/ThemeToggle';
import { LayoutDashboard, ClipboardList, Megaphone, LogOut, Menu } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTrigger, SheetTitle } from '@/components/ui/sheet';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useState } from 'react';

const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/logger', label: 'Logger', icon: ClipboardList },
  { to: '/marketing', label: 'Marketing', icon: Megaphone },
];

function NavItem({ to, label, icon: Icon, onClick }) {
  return (
    <NavLink
      to={to}
      onClick={onClick}
      data-testid={`nav-${label.toLowerCase()}`}
      className={({ isActive }) =>
        `flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all ${
          isActive
            ? 'bg-primary text-primary-foreground shadow-md'
            : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
        }`
      }
    >
      <Icon className="h-5 w-5" />
      <span>{label}</span>
    </NavLink>
  );
}

export default function AppLayout({ children }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/auth');
  };

  const initials = user?.name?.split(' ').map(n => n[0]).join('').toUpperCase() || '?';

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex md:w-72 flex-col border-r bg-card p-6">
        <div className="mb-8">
          <h1 className="text-2xl font-bold tracking-tight text-primary heading-font">
            SheLeads
          </h1>
          <p className="text-xs text-muted-foreground mt-1 uppercase tracking-widest">Empowering Women</p>
        </div>

        <nav className="flex-1 space-y-1">
          {navItems.map(item => (
            <NavItem key={item.to} {...item} />
          ))}
        </nav>

        <div className="mt-auto space-y-4">
          <ThemeToggle />
          <div className="flex items-center gap-3 p-3 rounded-xl bg-muted/30">
            <Avatar className="h-9 w-9">
              <AvatarFallback className="bg-primary text-primary-foreground text-sm">{initials}</AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium truncate">{user?.name || 'User'}</p>
              <p className="text-xs text-muted-foreground truncate">{user?.email || ''}</p>
            </div>
            <Button variant="ghost" size="icon" onClick={handleLogout} data-testid="logout-btn" className="h-8 w-8">
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-h-0">
        {/* Mobile Header */}
        <header className="md:hidden flex items-center justify-between px-4 py-3 border-b bg-card">
          <h1 className="text-lg font-bold text-primary heading-font">SheLeads</h1>
          <div className="flex items-center gap-2">
            <ThemeToggle />
            <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" data-testid="mobile-menu-btn">
                  <Menu className="h-5 w-5" />
                </Button>
              </SheetTrigger>
              <SheetContent side="left" className="w-72">
                <SheetTitle className="text-primary heading-font">SheLeads</SheetTitle>
                <nav className="mt-8 space-y-1">
                  {navItems.map(item => (
                    <NavItem key={item.to} {...item} onClick={() => setMobileOpen(false)} />
                  ))}
                </nav>
                <div className="mt-8">
                  <Button variant="ghost" onClick={handleLogout} className="w-full justify-start gap-3" data-testid="mobile-logout-btn">
                    <LogOut className="h-4 w-4" /> Logout
                  </Button>
                </div>
              </SheetContent>
            </Sheet>
          </div>
        </header>

        {/* Mobile Bottom Nav */}
        <nav className="md:hidden fixed bottom-0 left-0 right-0 z-50 flex items-center justify-around border-t bg-card py-2">
          {navItems.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              data-testid={`mobile-nav-${label.toLowerCase()}`}
              className={({ isActive }) =>
                `flex flex-col items-center gap-1 px-3 py-1 rounded-lg text-xs transition-colors ${
                  isActive ? 'text-primary' : 'text-muted-foreground'
                }`
              }
            >
              <Icon className="h-5 w-5" />
              <span>{label}</span>
            </NavLink>
          ))}
        </nav>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-4 md:p-8 pb-20 md:pb-8">
          <div className="animate-fade-in">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
