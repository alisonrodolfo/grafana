import React from 'react';
import { DataQuery } from '@grafana/data';
import { Button, Field, Input, InputControl, Modal } from '@grafana/ui';
import { FolderPicker } from 'app/core/components/Select/FolderPicker';
import { useForm } from 'react-hook-form';
import { SaveToNewDashboardDTO } from './addToDashboard';

type FormDTO = SaveToNewDashboardDTO;

interface Props {
  onClose: () => void;
  queries: DataQuery[];
  visualization: string;
  onSave: (data: FormDTO, redirect: boolean) => Promise<void | { message: string; status: string }>;
}

function withRedirect<T extends any[]>(fn: (redirect: boolean, ...args: T) => {}, redirect: boolean) {
  return async (...args: T) => fn(redirect, ...args);
}

export const AddToDashboardModal = ({ onClose, queries, visualization, onSave }: Props) => {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<FormDTO, { error: string }>({ defaultValues: { queries, visualization } });

  const onSubmit = async (withRedirect: boolean, data: FormDTO) => {
    const error = await onSave(data, withRedirect);

    if (error) {
      switch (error.status) {
        case 'name-exists':
        case 'empty-name':
        case 'name-match':
          setError('dashboardName', { message: error.message, shouldFocus: true });
          break;
        default:
        // TODO: Other unknown errors may happen, we should handle them by displaying an error message
      }
    }
  };

  return (
    <Modal
      // TODO: we can add multiple queries, shall we change the title?
      title="Add query to dashboard"
      onDismiss={onClose}
      isOpen
    >
      <form>
        <input type="hidden" {...register('queries')} />
        <input type="hidden" {...register('visualization')} />

        <Field label="Dashboard name" error={errors.dashboardName?.message} invalid={!!errors.dashboardName}>
          <Input
            id="dahboard_name"
            {...register('dashboardName', {
              shouldUnregister: true,
              required: { value: true, message: 'This field is required' },
            })}
            // we set default value here instead of in useForm because this input will be unregistered when switching
            // to "Existing Dashboard" and default values are not populated with manually registered
            // inputs (ie. when switching back to "New Dashboard")
            defaultValue="New dashboard (Explore)"
          />
        </Field>

        <Field label="Folder" error={errors.folderId?.message} invalid={!!errors.folderId}>
          <InputControl
            render={({ field: { ref, onChange, ...field } }) => (
              <FolderPicker onChange={(e) => onChange(e.id)} {...field} enableCreateNew inputId="folder" />
            )}
            control={control}
            name="folderId"
            shouldUnregister
            rules={{ required: { value: true, message: 'Select a valid folder to save your dashboard in' } }}
          />
        </Field>

        <Modal.ButtonRow>
          <Button type="reset" onClick={onClose} fill="outline" variant="secondary" disabled={isSubmitting}>
            Cancel
          </Button>
          <Button
            type="submit"
            onClick={handleSubmit(withRedirect(onSubmit, false))}
            variant="secondary"
            icon="compass"
            disabled={isSubmitting}
          >
            Save and keep exploring
          </Button>
          <Button
            type="submit"
            onClick={handleSubmit(withRedirect(onSubmit, true))}
            variant="primary"
            icon="save"
            disabled={isSubmitting}
          >
            Save and go to dashboard
          </Button>
        </Modal.ButtonRow>
      </form>
    </Modal>
  );
};
