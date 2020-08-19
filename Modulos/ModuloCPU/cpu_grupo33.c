#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>

#include <linux/sched.h>        // for_each_process, pr_info
#include <linux/sched/signal.h> 

#define BUFSIZE  	150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribir informacion del cpu.");
MODULE_AUTHOR("Jeannira Del Rosario Sic Menéndez - 201602434\nFernando Vidal Ruiz Piox -  201503984");

struct sysinfo inf;

int cont = 0;
long total_memoria = 0;
void dfs(struct task_struct *task, struct seq_file * m, int num)
{
  struct task_struct *task_next;
  struct list_head *list;
  
  int tmp = cont;
  cont = 1;
  
  int cant = 0;
  list_for_each(list, &task->children) {
    task_next = list_entry(list, struct task_struct, sibling);
    if(task_next == 0){
        continue;
    }
    cant = cant+1;
  }
  
  list_for_each(list, &task->children) {
    task_next = list_entry(list, struct task_struct, sibling);

    if(task_next == 0){
        continue;
    }
    seq_printf(m,"      {\n");
    //============= PID =====================
    seq_printf(m,"      \"PID\":%d,\n",task_next->pid);

	//============ NOMBRE ===================
    seq_printf(m,"      \"Nombre\":\"%s\",\n",task_next->comm);

    //============ USUARIO ==================
    seq_printf(m,"      \"Usuario\":\"%d\",\n",task_next->cred->uid);
    
    //============ ESTADO ===================
    if(task_next->state == -1){
        seq_printf(m,"      \"Estado\":\"NOT EXECUTABLE\"");
    }else if(task_next->state == 0){
        seq_printf(m,"      \"Estado\":\"RUNNING\"");
    }else if(task_next->state == 1){
        seq_printf(m,"      \"Estado\":\"INTERRUPTIBLE\"");
    }else if(task_next->state == 2){
        seq_printf(m,"      \"Estado\":\"UNINTERRUPTIBLE\"");
    }else if(task_next->state == 4){ //STOPPED
        seq_printf(m,"      \"Estado\":\"STOPPED\"");
    }else if(task_next->state == 8){ //TRACED
        seq_printf(m,"      \"Estado\":\"TRACED\"");
    }else if(task_next->state == 16){
        seq_printf(m,"      \"Estado\":\"ZOMBIE\"");
    }else{
        seq_printf(m,"      \"Estado\":\"EXCLUSIVE\"");
    }
    seq_printf(m,",\n");

    //============= Porcentaje RAM ============
    long memoria_utilizada = 0;
    if(task_next->mm){
        memoria_utilizada = (task_next->mm->total_vm <<(PAGE_SHIFT -10));
        memoria_utilizada = (memoria_utilizada/58603);
    }

    //memoria_utilizada*100/total_memoria
    seq_printf(m,"      \"PorcentajeRam\": %i",memoria_utilizada);
    seq_printf(m,",\n");

    //============= PPID ======================
    seq_printf(m,"      \"PPID\":%d\n",num);
    
    if(cant > 1){
        seq_printf(m,"      },\n");
        cant = cant - 1;
    }else {
        if(num == 2){
            seq_printf(m,"      }\n");    
        }else{
            seq_printf(m,"      },\n");
        }
    }
    dfs(task_next, m, task_next->pid); //hijos
  }
  cont = tmp;
}

static int escribir_archivo(struct seq_file * m, void *v) {	
    cont = 1;
    si_meminfo(&inf);
    total_memoria 	= (inf.totalram * 4)/1024;
    seq_printf(m,"{\n");
    seq_printf(m,"      \"struct_lista_procesos\":[\n");
    dfs(&init_task, m, 0);
    seq_printf(m,"      ]\n");
    seq_printf(m,"}\n");
    return 0;
}

static int al_abrir(struct inode *inode, struct  file *file) {
  return single_open(file, escribir_archivo, NULL);
}

static struct file_operations operaciones =
{    
    .open = al_abrir,
    .read = seq_read
};

static int iniciar(void)
{
    proc_create("cpu_grupo33", 0, NULL, &operaciones);
    printk(KERN_INFO "Hola Procesos!\nnombre1:Jeannira Del Rosario Sic Menéndez\nnombre2: Fernando Vidal Ruiz Piox\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("cpu_grupo33", NULL);
    printk(KERN_INFO "Adios Procesos!\nCurso: Sistemas Operativos 2\n");
}
 
module_init(iniciar);
module_exit(salir); 