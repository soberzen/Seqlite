import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@qingmo/shadcn/components/button';
import { EmojiPicker } from '@qingmo/shadcn/components/emoji-picker';
import { Field, FieldError } from '@qingmo/shadcn/components/field';
import { Input } from '@qingmo/shadcn/components/input';
import { Label } from '@qingmo/shadcn/components/label';
import { Textarea } from '@qingmo/shadcn/components/textarea';
import { FolderClosedIcon } from 'lucide-react';
import {
  useCallback,
  useEffect,
  useId,
  useMemo,
  useRef,
  type RefObject,
} from 'react';
import { useController, useForm, type Control } from 'react-hook-form';

import { Model } from '@/components/ui/model';
import {
  CreateSpaceSchema,
  spaceDescriptionLimit,
  spaceNameLimit,
  type CreateSpaceDto,
  type Space,
} from '@/service/api/space';

export type SpaceDialogRequest =
  | {
      mode: 'create';
      initialValues?: Partial<CreateSpaceDto>;
    }
  | { mode: 'update'; spaceId: Space['id']; initialValues: CreateSpaceDto };

type CreateSpaceDialogProps = {
  open: boolean;
  request: SpaceDialogRequest;
  onOpenChange: (open: boolean) => void;
};
type FieldProps = {
  control: Control<CreateSpaceDto>;
};

function SpaceNameField({
  control,
  inputRef,
}: FieldProps & { inputRef: RefObject<HTMLInputElement | null> }) {
  const id = useId();
  const countId = `${id}-count`;
  const errorId = `${id}-error`;

  const {
    field: { ref: nameFieldRef, ...nameField },
    fieldState: { error },
  } = useController({ name: 'name', control });

  const { field: emojiField } = useController({ name: 'emoji', control });

  const setInputRef = useCallback(
    (input: HTMLInputElement | null) => {
      inputRef.current = input;
      nameFieldRef(input);
    },
    [inputRef, nameFieldRef],
  );
  return (
    <Field
      className='gap-1.5'
      data-invalid={Boolean(error)}
    >
      <div className='relative'>
        <EmojiPicker
          modal
          value={emojiField.value}
          onSelect={(emoji) => emojiField.onChange(emoji)}
          onClear={() => emojiField.onChange(null)}
        >
          <Button
            type='button'
            variant='ghost'
            size='icon-xs'
            aria-label='选择空间图标'
            className='absolute inset-y-px left-px h-auto w-11 rounded-l-[calc(var(--radius-lg)-1px)] rounded-r-none border-0 border-r border-input text-muted-foreground'
          >
            {emojiField.value ?? (
              <FolderClosedIcon
                aria-hidden
                className='size-4'
              />
            )}
          </Button>
        </EmojiPicker>
        <Label
          htmlFor={id}
          className='sr-only'
        >
          空间名称
        </Label>
        <Input
          {...nameField}
          id={id}
          ref={setInputRef}
          type='text'
          className='pr-16 pl-14 h-10'
          required
          maxLength={spaceNameLimit}
          placeholder='空间名称'
          aria-invalid={Boolean(error)}
          aria-describedby={error ? `${countId} ${errorId}` : countId}
        />
        <output
          id={countId}
          htmlFor={id}
          aria-live='off'
          className='pointer-events-none absolute right-2.5 bottom-1.5 text-[11px] leading-4 text-muted-foreground tabular-nums'
        >
          {nameField.value.length} / {spaceNameLimit}
        </output>
      </div>
      {error && (
        <FieldError
          id={errorId}
          errors={[error]}
        />
      )}
    </Field>
  );
}

function SpaceDescriptionFieldLabel({ id }: { id: string }) {
  return (
    <div className='flex items-center justify-between'>
      <Label
        htmlFor={id}
        className='text-dialog-foreground'
      >
        描述
      </Label>
      <span className='text-xs text-dialog-muted'>选填</span>
    </div>
  );
}

function SpaceDescriptionField({ control }: FieldProps) {
  const id = useId();
  const countId = `${id}-count`;
  const errorId = `${id}-error`;

  const {
    field,
    fieldState: { error },
  } = useController({ name: 'description', control });

  return (
    <Field data-invalid={Boolean(errorId)}>
      <SpaceDescriptionFieldLabel id={id} />
      <div
        data-invalid={Boolean(error)}
        className='grid min-h-24 resize-y grid-rows-[1fr_auto] overflow-hidden rounded-lg border border-input transition-colors focus-within:border-ring focus-within:ring-3 focus-within:ring-ring/50 data-[invalid=true]:border-destructive data-[invalid=true]:ring-3 data-[invalid=true]:ring-destructive/20 dark:bg-input/30'
      >
        <Textarea
          {...field}
          id={id}
          rows={2}
          className='field-sizing-fixed min-h-0 resize-none rounded-b-none border-0 pb-0 text-dialog-foreground focus-visible:ring-0 aria-invalid:ring-0 dark:bg-transparent'
          placeholder='添加空间描述，让成员了解这里的内容…'
          maxLength={spaceDescriptionLimit}
          aria-invalid={Boolean(errorId)}
          aria-describedby={error ? `${countId} ${errorId}` : countId}
        />
        <output
          id={countId}
          htmlFor={id}
          aria-live='off'
          className='pr-6 pb-2 pl-2.5 text-right text-[11px] leading-4 text-muted-foreground tabular-nums'
        >
          {field.value.length} / {spaceDescriptionLimit}
        </output>
      </div>
      {error && (
        <FieldError
          id={errorId}
          errors={[error]}
        />
      )}
    </Field>
  );
}

export function CreateSpaceDialog({
  open,
  request,
  onOpenChange,
}: CreateSpaceDialogProps) {
  const formId = useId();
  const nameInputRef = useRef<HTMLInputElement>(null);
  const { control, handleSubmit, reset } = useForm({
    resolver: zodResolver(CreateSpaceSchema),
    defaultValues: { name: '', emoji: null, description: '' },
  });

  const isUpdating = useMemo(() => request.mode === 'update', [request.mode]);

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) {
      reset();
    }
    onOpenChange(nextOpen);
  }

  useEffect(() => {
    if (!open) return;

    reset({
      name: request.initialValues?.name ?? '',
      emoji: request.initialValues?.emoji ?? null,
      description: request.initialValues?.description ?? '',
    });
  }, [open, request, reset]);

  return (
    <Model
      open={open}
      onOpenChange={handleOpenChange}
      title='创建空间'
      initialFocusRef={nameInputRef}
      contentClassName='md:max-w-140'
      footer={
        <>
          <Button
            type='button'
            variant='ghost'
            className='px-4 text-dialog-secondary'
            onClick={() => handleOpenChange(false)}
          >
            取消
          </Button>
          <Button
            type='submit'
            form={formId}
            className='px-4'
          >
            {isUpdating ? '保存修改' : '创建空间'}
          </Button>
        </>
      }
    >
      <form
        id={formId}
        className='grid gap-4'
        noValidate
        onSubmit={handleSubmit(() => undefined)}
      >
        <SpaceNameField
          control={control}
          inputRef={nameInputRef}
        />
        <SpaceDescriptionField control={control} />
      </form>
    </Model>
  );
}
