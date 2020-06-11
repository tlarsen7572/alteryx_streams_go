long __declspec(dllexport) Controller(int nToolID,
   	void * pXmlProperties,
   	void *pEngineInterface,
   	void *r_pluginInterface);

long __declspec(dllexport) Interval(int nToolID,
   	void * pXmlProperties,
   	void *pEngineInterface,
   	void *r_pluginInterface);

long __declspec(dllexport) CombineLatest(int nToolID,
   	void * pXmlProperties,
   	void *pEngineInterface,
   	void *r_pluginInterface);
