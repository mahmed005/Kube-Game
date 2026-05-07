import './globals.css';

export default function Layout({children}: {children: React.ReactNode}) {
  return <html className="h-screen w-screen">
    <body className="h-full w-full">
      <header className="h-[100px] w-full bg-black/8 p-10 flex items-center justify-between">
        <h2 className='text-[26px] md:text-[40px] font-medium'>Kube Game</h2>
        <button className="bg-warning text-[18px] leading-loose cursor-pointer">Leave Game</button>
      </header>
      <main>
        {children}
      </main>
    </body>
  </html>
}