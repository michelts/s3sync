import '../styles/globals.css';

function MyApp({ Component, pageProps }): React.FC {
  return (
    <>
      <div className="w-full bg-blue-900">
        Test
      </div>
      <div className="container w-full md:w-96 mx-auto p-3">
        <Component {...pageProps} />
      </div>
    </>
  )
}

export default MyApp
