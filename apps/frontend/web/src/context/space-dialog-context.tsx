import {
  createContext,
  useCallback,
  useContext,
  useState,
  type PropsWithChildren,
} from 'react';

import {
  CreateSpaceDialog,
  type SpaceDialogRequest,
} from '@/views/space/components/create-space-dialog';

type OpenSpaceDialog = (request: SpaceDialogRequest) => void;

const SpaceDialogContext = createContext<OpenSpaceDialog | null>(null);

export function SpaceDialogProvider({ children }: PropsWithChildren) {
  const [dialog, setDialog] = useState<{
    open: boolean;
    request: SpaceDialogRequest;
  }>({
    open: false,
    request: { mode: 'create' },
  });

  const openSpaceDialog = useCallback<OpenSpaceDialog>((request) => {
    setDialog({ open: true, request: { ...request } });
  }, []);

  return (
    <SpaceDialogContext.Provider value={openSpaceDialog}>
      {children}
      <CreateSpaceDialog
        open={dialog.open}
        request={dialog.request}
        onOpenChange={(open) => {
          setDialog((current) => ({ ...current, open }));
        }}
      />
    </SpaceDialogContext.Provider>
  );
}

export function useOpenSpaceDialog() {
  const openSpaceDialog = useContext(SpaceDialogContext);
  if (!openSpaceDialog) {
    throw new Error(
      'useOpenSpaceDialog must be used within SpaceDialogProvider',
    );
  }
  return openSpaceDialog;
}
